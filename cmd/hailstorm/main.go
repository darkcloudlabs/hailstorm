package main

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/cluster"
	"github.com/anthdm/hollywood/remote"
	"github.com/darkcloudlabs/hailstorm/pkg/api"
)

func main() {
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

	c.Engine().Spawn(api.NewServer(c), "api")
	select {}
}
