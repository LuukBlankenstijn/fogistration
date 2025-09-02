package main

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func main() {
	ctx := context.Background()
	var cfg config.ClientConfig
	err := config.LoadFlags(&cfg)
	if err != nil {
		logging.Fatal("failed to get config", err)
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		logging.Fatal("Failed to create client", err)
	}

	c.RegisterHandler(&client.UpdateHandler{})

	if err = c.StartReceiving(ctx); err != nil {
		logging.Fatal("failed receiving messages", err)
	}

}
