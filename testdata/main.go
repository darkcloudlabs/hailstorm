package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darkcloudlabs/hailstorm/hailstorm"
)

func init() {
	hailstorm.Handle(myHandler)
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	fmt.Println("body: ", string(body))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("test", "test")
	json.NewEncoder(w).Encode(map[string]string{"message": "hello world"})

}

func main() {}
