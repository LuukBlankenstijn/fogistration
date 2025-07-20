package wrapper

import (
	"context"
	"fmt"

	gen "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	"slices"
)

// ListContests fetches the contest list.
// params may be nil. reqEditors are passed through to the generated client.
func (c *Client) ListContests(
	ctx context.Context,
	params *gen.GetV4AppApiContestListParams,
	reqEditors ...gen.RequestEditorFn,
) ([]gen.Contest, error) {

	resp, err := c.raw.GetV4AppApiContestListWithResponse(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	// Success
	if resp.JSON200 != nil {
		// resp.JSON200 is *[]gen.Contest
		return slices.Clone((*resp.JSON200)), nil
	}

	// Known error payloads
	if resp.JSON400 != nil {
		return nil, fmt.Errorf("cmdhandler: invalid request: %v", *resp.JSON400)
	}
	if resp.JSON403 != nil {
		return nil, fmt.Errorf("cmdhandler: unauthorized: %v", *resp.JSON403)
	}
	if resp.JSON404 != nil {
		return nil, fmt.Errorf("cmdhandler: not found: %v", *resp.JSON404)
	}

	// Fallback: unexpected status/content
	return nil, apiResponseErr("GetV4AppApiContestList", resp.HTTPResponse, resp.Body)
}
