package wrapper

import (
	"context"
	"fmt"
	"slices"

	gen "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
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
func (c *Client) UpdateUser(
	ctx context.Context,
	id string,
	params gen.UpdateUser,
	reqEditors ...gen.RequestEditorFn,
) (*gen.User, error) {

	resp, err := c.raw.PatchV4AppApiUserUpdateWithResponse(
		ctx,
		id,
		params,
		reqEditors...,
	)
	if err != nil {
		return nil, err
	}

	if resp.JSON201 != nil {
		return resp.JSON201, nil
	}

	// Known error payloads
	if resp.JSON400 != nil {

		logging.Info("Raw error: %s\n", string(resp.Body))
		return nil, fmt.Errorf("cmdhandler: invalid request: %v", *resp.JSON400)
	}
	if resp.JSON403 != nil {
		return nil, fmt.Errorf("cmdhandler: unauthorized: %v", *resp.JSON403)
	}
	if resp.JSON404 != nil {
		return nil, fmt.Errorf("cmdhandler: not found: %v", *resp.JSON404)
	}

	return nil, apiResponseErr("PutV4AppApiUserUpdateWithBodyWithResponse", resp.HTTPResponse, resp.Body)
}
