package main

import (
	"io"
	"net/http"

	hailstorm "github.com/darkcloudlabs/hailstorm/sdk"
)

func init() {
	hailstorm.Handle(myHandler)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://google.com")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	io.Copy(w, resp.Body)
}

func main() {}
