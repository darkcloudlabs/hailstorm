package hailstorm

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"unsafe"
)

var handler http.HandlerFunc

func Handle(h http.HandlerFunc) {
	handler = h
}

var buffers = make(map[uintptr][]byte, 0)

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
	// TODO: handle this buff.
	buf := readBufferFromMemory(&ptr, n)
	_ = buf

	req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("foo")))
	rw := &ResponseWriter{}

	handler(rw, req)

	b := rw.buffer.Bytes()
	respb := make([]byte, 4+len(b))
	binary.LittleEndian.PutUint32(respb, uint32(len(b)))

	copy(respb[4:], b)
	p := makeBuffer(respb)

	return p
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
