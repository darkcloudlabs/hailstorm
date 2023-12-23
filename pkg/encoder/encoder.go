package encoder

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net/http"
)

var (
	errorBody    = []byte("Internal Server Error\n")
	errorHeaders = []byte("{}")
)

func SafeEncode(body []byte, headers []byte, statusCode int) []byte {
	// [body length][headers length][body][headers][status code]
	respb := make([]byte, 4+4+len(body)+len(headers)+4)
	binary.LittleEndian.PutUint32(respb, uint32(len(body)))
	binary.LittleEndian.PutUint32(respb[4:], uint32(len(headers)))
	copy(respb[8:], body)
	copy(respb[8+len(body):], headers)
	binary.LittleEndian.PutUint32(respb[8+len(body)+len(headers):], uint32(statusCode))
	return respb
}

func SafeDecode(res []byte) (body []byte, headers []byte, statusCode int) {
	// handle if the response returned 0
	if len(res) < 8 {
		return errorBody, errorHeaders, http.StatusInternalServerError
	}

	bodyLength := binary.LittleEndian.Uint32(res)
	headersLength := binary.LittleEndian.Uint32(res[4:])
	body = res[8 : 8+bodyLength]
	headers = res[8+bodyLength : 8+bodyLength+headersLength]
	statusCode = int(binary.LittleEndian.Uint32(res[8+bodyLength+headersLength:]))
	return body, headers, statusCode
}

func EncodeHeaders(h http.Header) []byte {
	if h == nil {
		return json.RawMessage("{}")
	}

	headers, _ := json.Marshal(h)
	return headers
}

func DecodeHeaders(headers []byte) http.Header {
	h := make(http.Header)
	_ = json.Unmarshal(headers, &h)
	return h
}

type Request struct {
	Body   []byte
	Method string
	URL    string
}

// using json for now, but we can use protobuf or other binary format
func EncodeRequest(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, nil
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req := Request{
		Method: r.Method,
		URL:    r.URL.Path,
		Body:   b,
	}

	return json.Marshal(req)
}

func DecodeRequest(b []byte) (*http.Request, error) {
	req := Request{}
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	r, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}

	return r, nil
}
