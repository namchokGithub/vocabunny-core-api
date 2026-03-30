package main

import (
	"context"
	"log"

	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/namchokGithub/vocabunny-core-api/protocol"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := protocol.Initialize(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("initialize app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("start app: %v", err)
	}
}
