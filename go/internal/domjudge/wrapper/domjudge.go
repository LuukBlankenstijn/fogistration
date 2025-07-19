package domjudge

import (
	"context"
	"fmt"
	"net/http"

	gen "github.com/LuukBlankenstijn/fogistration/internal/domjudge/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
)

type Client struct {
	cfg config.DomJudgeConfig
	hc  *http.Client
	raw *gen.ClientWithResponses
}

func NewClient(ctx context.Context, cfg config.DomJudgeConfig) (*Client, error) {
	hc := &http.Client{
		Transport: gen.NewBasicAuthRT(cfg.Username, cfg.Password, http.DefaultTransport),
	}

	raw, err := gen.NewClientWithResponses(cfg.DJHost, gen.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}
	return &Client{
		cfg: cfg,
		hc:  hc,
		raw: raw,
	}, nil
}

// func (c *Client) Raw() *gen.ClientWithResponses {
// 	return c.raw
// }

func apiResponseErr(op string, r *http.Response, body []byte) error {
	if r == nil {
		return fmt.Errorf("%s: no http response", op)
	}
	return fmt.Errorf("%s: unexpected status %d: %s", op, r.StatusCode, string(body))
}
