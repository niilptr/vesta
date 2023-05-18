package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"vesta/x/twin/types"
)

func (k msgServer) Train(goCtx context.Context, msg *types.MsgTrain) (*types.MsgTrainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgTrainResponse{}, nil
}
