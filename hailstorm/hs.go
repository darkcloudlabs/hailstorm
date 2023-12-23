package hailstorm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/darkcloudlabs/hailstorm/pkg/encoder"
)

var handlers = make(map[uint32]http.HandlerFunc, 0)

func Handle(h http.HandlerFunc) {
	handlers[functionCount()] = h
}

//export function_count
func functionCount() uint32 {
	return uint32(len(handlers))
}

var mem = make(map[uint32][]byte, 0)

func makeBuffer(b []byte) (ptr uint32) {
	ptr = uint32(uintptr(unsafe.Pointer(&b[0])))
	mem[ptr] = b
	return ptr
}

func readBufferFromMemory(ptr, size uint32) ([]byte, error) {
	b, ok := mem[ptr]
	if !ok {
		return nil, fmt.Errorf("pointer %d not found in memory", ptr)
	}
	return b[:size], nil
}

//export alloc
func alloc(size uint32) (loc uint32) {
	return makeBuffer(make([]byte, size))
}

//export handle_http_request
func handleHandleHTTPRequest(ptr uint32, size uint32, id uint32) uint32 {
	buf, err := readBufferFromMemory(ptr, size)
	if err != nil {
		fmt.Printf("hs.go: error reading buffer from memory: %s\n", err)
		return 0
	}

	req, err := encoder.DecodeRequest(buf)
	if err != nil {
		fmt.Printf("hs.go: error decoding request: %s\n", err)
		return 0
	}

	rw := &ResponseWriter{header: make(http.Header)}

	handler := handlers[id]
	handler(rw, req)

	bodyBytes := rw.buffer.Bytes()
	headersBytes := encoder.EncodeHeaders(rw.header)

	b := encoder.SafeEncode(bodyBytes, headersBytes, rw.statusCode)

	respb := make([]byte, 4+len(b))
	binary.LittleEndian.PutUint32(respb, uint32(len(b)))

	copy(respb[4:], b)
	p := makeBuffer(respb)

	return p
}

type ResponseWriter struct {
	buffer     bytes.Buffer
	header     http.Header
	statusCode int
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.statusCode = status
}
