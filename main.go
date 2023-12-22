package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/vmihailenco/msgpack/v5"
)

// type Response struct{

// }

type Request struct {
	Method string
	URL    string
	Body   []byte
}

func main() {
	// HandlerOnTheHost(Response{}, Request{})

	http.HandleFunc("/", InternalHandler)

	http.ListenAndServe(":3000", nil)
}

func InternalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new incoming request on the host")

	b, err := os.ReadFile("app.wasm")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithDebugInfoEnabled(true))
	defer runtime.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	mod, err := runtime.InstantiateWithConfig(ctx, b, wazero.NewModuleConfig().WithStdout(os.Stdout))
	if err != nil {
		log.Fatal(err)
	}

	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	req := Request{
		Method: r.Method,
		URL:    r.URL.Path,
		Body:   rbody,
	}

	reqb, err := msgpack.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := mod.ExportedFunction("alloc").Call(ctx, uint64(len(reqb)))
	if err != nil {
		log.Fatal(err)
	}

	mod.Memory().Write(uint32(res[0]), reqb)

	res, err = mod.ExportedFunction("handle_http_request").Call(ctx, res[0], uint64(len(reqb)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	_ = mod
	w.Write([]byte("good"))
}

// go:embed main.wasm
var wasmbytes []byte

func bar() uint32 {
	fmt.Println("from the host")
	return 88
}

func makeRuntime() {

	// mod, err := runtime.NewHostModuleBuilder("foo").
	// 	NewFunctionBuilder().
	// 	WithGoModuleFunction()
	// 	WithFunc(bar).
	// 	Export("bar").
	// 	Instantiate(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//_ = mod
	// ok := mod.Memory().WriteString(10, "hello world")
	// if !ok {
	// 	log.Fatal("out of range")
	// }

}
