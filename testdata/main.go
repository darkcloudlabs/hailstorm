package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darkcloudlabs/hailstorm/hailstorm"
)

func init() {
	hailstorm.Handle(handler1)
	hailstorm.Handle(handler2)
	hailstorm.Handle(handler3)
}

func handler1(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	fmt.Println("body: ", string(body))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"function": "handler1"})
}

func handler2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"function": "handler2"})
}

func handler3(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"function": "handler3"})
}

func main() {}
