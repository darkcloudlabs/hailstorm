package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/cluster"
	"github.com/anthdm/hollywood/remote"
	"github.com/darkcloudlabs/hailstorm/pkg/api"
	"github.com/darkcloudlabs/hailstorm/pkg/proxy"
	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/google/uuid"
)

func main() {
	memstore := store.NewMemoryStore()

	r := remote.New(remote.Config{
		ListenAddr: "127.0.0.1:30000",
	})
	e, err := actor.NewEngine(actor.EngineOptRemote(r))
	if err != nil {
		log.Fatal(err)
	}
	c, err := cluster.New(cluster.Config{
		ClusterProvider: cluster.NewSelfManagedProvider(),
		ID:              "server",
		Engine:          e,
		Region:          "eu-west",
	})
	c.Start()

	c.Engine().Spawn(api.NewServer(c, memstore), "api")

	seed(memstore)

	proxy := proxy.New(memstore)
	log.Fatal(proxy.Listen(":5000"))
}

func seed(store store.Store) {
	b, err := os.ReadFile("examples/go/app.wasm")
	if err != nil {
		log.Fatal(err)
	}
	app := types.App{
		ID:        uuid.New(),
		Name:      "My first Hailstorm app",
		CreatedAT: time.Now(),
	}
	store.CreateApp(&app)
	deploy := types.Deploy{
		ID:        uuid.New(),
		AppID:     app.ID,
		Blob:      b,
		CreatedAT: time.Now(),
	}
	store.CreateDeploy(&deploy)
	fmt.Printf("My first Hailstorm app available localhost:5000/%s\n", deploy.ID)
}
