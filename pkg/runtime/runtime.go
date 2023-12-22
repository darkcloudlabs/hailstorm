package runtime

import (
	"github.com/anthdm/hollywood/actor"
)

type Runtime struct {
}

func New() actor.Producer {
	return func() actor.Receiver {
		return &Runtime{}
	}
}

func (r *Runtime) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		// TODO: very poccy, we need to spawn another runner, which we can send messages to
		// (stop, ...)
		go r.start()
	default:
		_ = msg
	}
}

func (r *Runtime) start() {
	// ctx := context.Background()
	// runtime := wazero.NewRuntime(ctx)
	// defer runtime.Close(ctx)

	// wasmModule, err := runtime.CompileModule(ctx, r.blob.Contents)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer wasmModule.Close(ctx)

	// builder := imports.NewBuilder().
	// 	WithName("foobarbaz").
	// 	WithSocketsExtension("auto", wasmModule)
	// 	// WithArgs(args...).
	// 	// WithEnv(envs...).
	// 	// WithDirs(dirs...).
	// 	// WithListens(listens...).
	// 	// WithDials(dials...).
	// 	// WithNonBlockingStdio(nonBlockingStdio).
	// 	// WithTracer(trace, os.Stderr, wasi.WithTracerStringSize(tracerStringSize)).
	// 	// WithMaxOpenFiles(maxOpenFiles).
	// 	// WithMaxOpenDirs(maxOpenDirs)

	// var system wasi.System
	// ctx, system, err = builder.Instantiate(ctx, runtime)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer system.Close(ctx)

	// // Probably only need to do this when the user give us the "http service" flag
	// wasi_http.DetectWasiHttp(wasmModule)

	// wasiHTTP := wasi_http.MakeWasiHTTP()
	// if err := wasiHTTP.Instantiate(ctx, runtime); err != nil {
	// 	log.Fatal(err)
	// }

	// instance, err := runtime.InstantiateModule(ctx, wasmModule, wazero.NewModuleConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// instance.Close(ctx)
}
