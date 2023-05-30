package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ConfirmBestTrainResultIs(goCtx context.Context, msg *types.MsgConfirmBestTrainResultIs) (*types.MsgConfirmBestTrainResultIsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: check signer authorization

	err := k.AddBestTrainResultToTrainingState(ctx, msg.Creator, msg.Hash)
	if err != nil {
		return &types.MsgConfirmBestTrainResultIsResponse{}, err
	}

	return &types.MsgConfirmBestTrainResultIsResponse{}, nil
}
