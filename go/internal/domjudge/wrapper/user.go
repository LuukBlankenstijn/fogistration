package domjudge

import (
	"context"
	"fmt"
	"slices"

	gen "github.com/LuukBlankenstijn/fogistration/internal/domjudge/client"
)

func (c *Client) ListUsers(
	ctx context.Context,
	params *gen.GetV4AppApiUserListParams,
	reqEditors ...gen.RequestEditorFn,
) ([]gen.User, error) {
	resp, err := c.raw.GetV4AppApiUserListWithResponse(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 != nil {
		// resp.JSON200 is *[]gen.Contest
		return slices.Clone((*resp.JSON200)), nil
	}

	// Known error payloads
	if resp.JSON400 != nil {
		return nil, fmt.Errorf("domjudge: invalid request: %v", *resp.JSON400)
	}
	if resp.JSON403 != nil {
		return nil, fmt.Errorf("domjudge: unauthorized: %v", *resp.JSON403)
	}
	if resp.JSON404 != nil {
		return nil, fmt.Errorf("domjudge: not found: %v", *resp.JSON404)
	}

	return nil, apiResponseErr("GetV4AppApiContestList", resp.HTTPResponse, resp.Body)
}
