package remotedialer

import (
	"os"
	"strconv"
)

var (
	readBufferSize, writeBufferSize, bufferSize, sessionReadBufferSize, sessionWirteBufferSize int
)

func init() {
	buffer := os.Getenv("TUNNEL_BUFFER_SIZE")
	bufferSize, _ = strconv.Atoi(buffer)
	if bufferSize == 0 {
		bufferSize = 4096
	}
	readBufferSize, writeBufferSize = GetBufferSetting()
	sessionReadBufferSize, sessionWirteBufferSize = GetSessionBuffer()
}

func GetBufferSetting() (int, int) {
	readBuffer := os.Getenv("TUNNEL_READ_BUFFER_SIZE")
	writeBuffer := os.Getenv("TUNNEL_WRITE_BUFFER_SIZE")
	readBufferSize, _ := strconv.Atoi(readBuffer)
	writeBufferSize, _ := strconv.Atoi(writeBuffer)
	if readBufferSize == 0 {
		readBufferSize = bufferSize
	}
	if writeBufferSize == 0 {
		writeBufferSize = bufferSize
	}
	return readBufferSize, writeBufferSize
}

func GetSessionBuffer() (int, int) {
	readBuffer := os.Getenv("SESSION_READ_BUFFER_SIZE")
	readBufferSize, _ := strconv.Atoi(readBuffer)
	if readBufferSize == 0 {
		readBufferSize = bufferSize
	}
	writeBuffer := os.Getenv("SESSION_WRITE_BUFFER_SIZE")
	writeBufferSize, _ := strconv.Atoi(writeBuffer)
	if writeBufferSize == 0 {
		writeBufferSize = bufferSize
	}
	return readBufferSize, writeBufferSize
}
