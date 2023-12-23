package runtime

import (
	"context"
	"net/http"
	"os"

	"github.com/darkcloudlabs/hailstorm/pkg/encoder"
	"github.com/tetratelabs/wazero"
	wapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

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

	reqb, err := encoder.EncodeRequest(r)
	if err != nil {
		return err
	}

	alloc := runtime.module.ExportedFunction("alloc")
	res, err := alloc.Call(ctx, uint64(len(reqb)))
	if err != nil {
		return err
	}

	ptr := uint32(res[0])
	runtime.module.Memory().Write(ptr, reqb)

	handleHTTP := runtime.module.ExportedFunction("handle_http_request")
	res, err = handleHTTP.Call(ctx, res[0], uint64(len(reqb)))
	if err != nil {
		return err
	}
	n, _ := runtime.module.Memory().ReadUint32Le(uint32(res[0]))
	resBytes, _ := runtime.module.Memory().Read(uint32(res[0]), n+4)

	bodyBytes, headersBytes, statusCode := encoder.SafeDecode(resBytes[4:])

	headers := encoder.DecodeHeaders(headersBytes)

	for k, v := range headers {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(bodyBytes)
	return err
}

func (runtime *Runtime) Close(ctx context.Context) error {
	return runtime.wruntime.Close(ctx)
}
