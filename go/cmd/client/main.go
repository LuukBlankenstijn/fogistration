package main

import (
	"context"
	"time"

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

	for {
		err := c.StartReceiving(ctx)
		if err != nil {
			logging.Error("failed receiving messages", err)
		}

		const (
			baseDelay = 500 * time.Millisecond
			maxDelay  = 10 * time.Second
		)
		for delay := baseDelay; err != nil; delay *= 2 {
			if delay > maxDelay {
				delay = maxDelay
			}
			select {
			case <-time.After(delay):
				err = c.StartReceiving(ctx)
				if err == nil {
					break
				}
				logging.Error("retry failed", err)
			case <-ctx.Done():
				return
			}
		}
	}
}
