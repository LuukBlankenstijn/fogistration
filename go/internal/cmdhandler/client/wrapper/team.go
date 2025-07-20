package wrapper

import (
	"context"
	"fmt"
	"slices"

	gen "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
)

func (c *Client) ListTeams(
	ctx context.Context,
	params *gen.GetV4AppApiTeamListParams,
	cid gen.Cid,
	reqEditors ...gen.RequestEditorFn,
) ([]gen.Team, error) {
	resp, err := c.raw.GetV4AppApiTeamListWithResponse(ctx, cid, params, reqEditors...)
	if err != nil {
		return nil, err
	}

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

	return nil, apiResponseErr("GetV4AppApiContestList", resp.HTTPResponse, resp.Body)
}
