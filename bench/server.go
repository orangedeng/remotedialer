package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

var (
	clients = map[string]*http.Client{}
	l       sync.Mutex
	counter int64
)

func authorizer(req *http.Request) (string, bool, error) {
	id := req.Header.Get("x-tunnel-id")
	return id, id != "", nil
}

func Client(server *remotedialer.Server, rw http.ResponseWriter, req *http.Request) {
	timeout := req.URL.Query().Get("timeout")
	if timeout == "" {
		timeout = "15"
	}

	vars := mux.Vars(req)
	clientKey := vars["id"]
	url := fmt.Sprintf("%s://%s%s?%s", vars["scheme"], vars["host"], vars["path"], req.URL.Query().Encode())
	client := getClient(server, clientKey, timeout)

	id := atomic.AddInt64(&counter, 1)
	logrus.Debugf("[%03d] REQ t=%s %s", id, timeout, url)

	resp, err := client.Get(url)
	if err != nil {
		logrus.Errorf("[%03d] REQ ERR t=%s %s: %v", id, timeout, url, err)
		remotedialer.DefaultErrorWriter(rw, req, 500, err)
		return
	}
	defer resp.Body.Close()

	logrus.Debugf("[%03d] REQ OK t=%s %s", id, timeout, url)
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
	logrus.Debugf("[%03d] REQ DONE t=%s %s", id, timeout, url)
}

func getClient(server *remotedialer.Server, clientKey, timeout string) *http.Client {
	l.Lock()
	defer l.Unlock()

	key := fmt.Sprintf("%s/%s", clientKey, timeout)
	client := clients[key]
	if client != nil {
		return client
	}

	dialer := server.Dialer(clientKey)
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: dialer,
		},
	}
	if timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err == nil {
			client.Timeout = time.Duration(t) * time.Second
		}
	}

	clients[key] = client
	return client
}
