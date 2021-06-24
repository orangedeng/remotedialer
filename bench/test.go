package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"github.com/tcnksm/go-httpstat"
)

var (
	listen     = flag.String("listen", ":8125", "Listen address")
	proxyAddr  = flag.String("proxy", ":8124", "Proxy server listen address")
	bufferSize = flag.Int64("buffer", 4096, "Buffer size of proxy server and client")
	debug      = flag.Bool("debug", false, "Enable test debug")
	testSizes  = []int64{
		4 * 1024,
		256 * 1024,
		1024 * 1024,
		4 * 1024 * 1024,
		16 * 1024 * 1024,
		256 * 1024 * 1024,
	}
)

func main() {
	flag.Parse()
	if *bufferSize > 0 {
		os.Setenv("BUFFER_SIZE", strconv.FormatInt(*bufferSize, 10))
		defer os.Unsetenv("BUFFER_SIZE")
	}
	// setup dummy http server
	go http.ListenAndServe(*listen, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		size := req.URL.Query().Get("size")
		filesize, _ := strconv.ParseInt(size, 10, 64)
		if filesize == 0 {
			filesize = 1024
		}
		rtnBytes := make([]byte, filesize)
		start := time.Now()
		rw.WriteHeader(200)
		if _, err := rw.Write(rtnBytes); err != nil {
			fmt.Println(err)
		}
		logrus.Debugf("request %s, %s", req.URL.String(), time.Since(start))
	}))
	if err := testDirect(testSizes...); err != nil {
		panic(err)
	}
	if err := setupProxy(); err != nil {
		panic(err)
	}
	if err := testProxy(testSizes...); err != nil {
		panic(err)
	}
}

func testDirect(sizes ...int64) error {
	return access("direct", "http://127.0.0.1"+*listen, sizes...)
}

func testProxy(sizes ...int64) error {
	return access("proxy", "http://127.0.0.1"+*proxyAddr+"/client/foo/http/127.0.0.1"+*listen+"/", sizes...)
}

func access(logPrefix, baseURL string, sizes ...int64) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	buffer := os.Getenv("BUFFER_SIZE")
	bufferSize, _ := strconv.Atoi(buffer)
	if bufferSize == 0 {
		bufferSize = 4096
	}

	for _, size := range sizes {
		query := u.Query()
		query.Set("size", strconv.FormatInt(size, 10))
		u.RawQuery = query.Encode()
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return err
		}
		var result httpstat.Result
		req = req.WithContext(httpstat.WithHTTPStat(req.Context(), &result))
		client := http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2: false,
				ReadBufferSize:    bufferSize,
				WriteBufferSize:   bufferSize,
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		}
		endTime := time.Now()
		result.End(endTime)
		printResult(logPrefix, size, endTime, &result)
	}

	return nil
}

func setupProxy() error {
	if err := setupProxyServer(); err != nil {
		return err
	}
	if err := setupProxyClient(); err != nil {
		return err
	}
	return nil
}

func setupProxyServer() error {
	var (
		peerID    = "server"
		peerToken = ""
	)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
		remotedialer.PrintTunnelData = true
	}

	handler := remotedialer.New(authorizer, remotedialer.DefaultErrorWriter)
	handler.PeerToken = peerToken
	handler.PeerID = peerID

	router := mux.NewRouter()
	router.Handle("/connect", handler)
	router.HandleFunc("/client/{id}/{scheme}/{host}{path:.*}", func(rw http.ResponseWriter, req *http.Request) {
		Client(handler, rw, req)
	})

	fmt.Println("Listening on ", *proxyAddr)
	go http.ListenAndServe(*proxyAddr, router)
	return nil
}

func setupProxyClient() error {
	addr := "ws://127.0.0.1" + *proxyAddr + "/connect"
	headers := http.Header{
		"X-Tunnel-ID": []string{"foo"},
	}
	connectedChan := make(chan int)
	go remotedialer.ClientConnect(context.Background(), addr, headers, nil, func(string, string) bool { return true }, func(ctx context.Context, session *remotedialer.Session) error {
		logrus.Infof("proxy client connected")
		close(connectedChan)
		return nil
	})
	<-connectedChan
	return nil
}

func printResult(prefix string, size int64, endTime time.Time, result *httpstat.Result) {
	logrus.Printf("mod: %s,\tsize: %d\tcontent transfer: %s\ttotal: %s", prefix, size, result.ContentTransfer(endTime), result.Total(endTime))
}
