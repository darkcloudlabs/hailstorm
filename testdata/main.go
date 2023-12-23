package main

import (
	"fmt"
	"net/http"

	"github.com/darkcloudlabs/hailstorm/hailstorm"
)

func init() {
	hailstorm.Handle(myHandler)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is from the guest application")

	w.Write([]byte("my response"))
}

func main() {}
