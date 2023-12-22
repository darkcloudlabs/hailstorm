package hailstorm

import (
	"bytes"
	"fmt"
	"net/http"
	"unsafe"
)

var handler http.HandlerFunc

func Handle(h http.HandlerFunc) {
	handler = h
}

var buffers = map[uintptr][]byte{}

func readBufferFromMemory(bufferPosition *uint32, length uint32) []byte {
	subjectBuffer := make([]byte, length)
	pointer := uintptr(unsafe.Pointer(bufferPosition))
	for i := 0; i < int(length); i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		subjectBuffer[i] = byte(s)
	}
	return subjectBuffer
}

func makeBuffer(buf []byte) uint32 {
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	buffers[unsafePtr] = buf
	return uint32(unsafePtr)
}

//export alloc
func alloc(size uint32) uint32 {
	return makeBuffer(make([]byte, size))
}

//export handle_http_request
func handleHandleHTTPRequest(ptr uint32, n uint32) uint32 {
	buf := readBufferFromMemory(&ptr, n)
	_ = buf

	fmt.Println(string(buffers[uintptr(ptr)]))

	req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("foo")))
	rw := &ResponseWriter{}

	handler(rw, req)

	b := rw.buffer.Bytes()
	p := makeBuffer(b)

	fmt.Println(string(b))
	fmt.Println(len(b))

	return p, uint32(len(b))
}

type ResponseWriter struct {
	buffer     bytes.Buffer
	statusCode int
}

func (w *ResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *ResponseWriter) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.statusCode = status
}
