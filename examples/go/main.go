package main

import (
	"net/http"
	"os"

	hailstorm "github.com/darkcloudlabs/hailstorm/sdk"
)

func init() {
	hailstorm.Handle(myHandler)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	foo := os.Getenv("FOO")

	resp := "this is the response: " + foo

	w.WriteHeader(500)
	w.Write([]byte(resp))
}

func main() {}
