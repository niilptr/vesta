package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"vesta/x/twin/types"
)

func (k msgServer) ConfirmBestTrainResultIs(goCtx context.Context, msg *types.MsgConfirmBestTrainResultIs) (*types.MsgConfirmBestTrainResultIsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgConfirmBestTrainResultIsResponse{}, nil
}
