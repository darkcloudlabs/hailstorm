package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/darkcloudlabs/hailstorm/pkg/api"
	"github.com/darkcloudlabs/hailstorm/pkg/proxy"
	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/google/uuid"
)

const (
	httpListenAddr = ":3000"
	httpProxyAddr  = ":5000"
)

func main() {
	memstore := store.NewMemoryStore()
	seed(memstore)
	server := api.NewServer(memstore)
	go func() {
		slog.Info("api server running", "port", httpListenAddr)
		log.Fatal(server.Listen(httpListenAddr))
	}()

	proxy := proxy.New(memstore)
	slog.Info("app proxy server running", "port", httpProxyAddr)
	log.Fatal(proxy.Listen(httpProxyAddr))

	log.Fatal(proxy.Listen(":5000"))
}

func seed(store store.Store) {
	b, err := os.ReadFile("examples/go/app.wasm")
	if err != nil {
		log.Fatal(err)
	}
	app := types.App{
		ID:          uuid.New(),
		Name:        "My first Hailstorm app",
		Environment: map[string]string{"FOO": "fooenv"},
		CreatedAT:   time.Now(),
	}
	store.CreateApp(&app)
	deploy := types.Deploy{
		ID:        uuid.MustParse("09248ef6-c401-4601-8928-5964d61f2c61"),
		AppID:     app.ID,
		Blob:      b,
		CreatedAT: time.Now(),
	}
	store.CreateDeploy(&deploy)
	fmt.Printf("My first Hailstorm app available localhost:5000/%s\n", deploy.ID)
}
