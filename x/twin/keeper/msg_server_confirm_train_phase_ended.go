package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"vesta/x/twin/types"
)

func (k msgServer) ConfirmTrainPhaseEnded(goCtx context.Context, msg *types.MsgConfirmTrainPhaseEnded) (*types.MsgConfirmTrainPhaseEndedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgConfirmTrainPhaseEndedResponse{}, nil
}
