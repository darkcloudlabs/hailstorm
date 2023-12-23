package runtime

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tetratelabs/wazero"
	wapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/vmihailenco/msgpack/v5"
)

type Request struct {
	Body   []byte
	Method string
	URL    string
}

type Runtime struct {
	wasmBlob    []byte
	compiledMod wazero.CompiledModule
	module      wapi.Module
	wruntime    wazero.Runtime
}

func New(blob []byte) (*Runtime, error) {
	var (
		ctx     = context.Background()
		config  = wazero.NewRuntimeConfig().WithDebugInfoEnabled(true)
		runtime = wazero.NewRuntimeWithConfig(ctx, config)
	)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	compiledMod, err := runtime.CompileModule(ctx, blob)
	if err != nil {
		return nil, err
	}

	modConfig := wazero.NewModuleConfig().WithStdout(os.Stdout)
	mod, err := runtime.InstantiateModule(ctx, compiledMod, modConfig)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		wasmBlob:    blob,
		compiledMod: compiledMod,
		module:      mod,
		wruntime:    runtime,
	}, nil
}

func (runtime *Runtime) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	req := Request{
		Method: r.Method,
		URL:    r.URL.Path,
		Body:   rbody,
	}

	reqb, err := msgpack.Marshal(req)
	if err != nil {
		return err
	}

	fmt.Println(reqb)

	alloc := runtime.module.ExportedFunction("alloc")
	res, err := alloc.Call(ctx, uint64(len(reqb)))
	if err != nil {
		return err
	}

	runtime.module.Memory().Write(uint32(res[0]), reqb)

	handleHTTP := runtime.module.ExportedFunction("handle_http_request")
	res, err = handleHTTP.Call(ctx, res[0], uint64(len(reqb)))
	if err != nil {
		return err
	}
	n, _ := runtime.module.Memory().ReadUint32Le(uint32(res[0]))
	resBytes, _ := runtime.module.Memory().Read(uint32(res[0]), n+4)

	fmt.Println(string(resBytes))
	_, err = w.Write(resBytes[4:])
	return err
}

func (runtime *Runtime) Close(ctx context.Context) error {
	return runtime.wruntime.Close(ctx)
}
