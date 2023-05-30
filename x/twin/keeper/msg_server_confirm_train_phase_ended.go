package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ConfirmTrainPhaseEnded(goCtx context.Context, msg *types.MsgConfirmTrainPhaseEnded) (*types.MsgConfirmTrainPhaseEndedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: check sender authorization

	err := k.AddTrainingPhaseEndedConfirmation(ctx, msg.Creator)
	if err != nil {
		return &types.MsgConfirmTrainPhaseEndedResponse{}, err
	}

	return &types.MsgConfirmTrainPhaseEndedResponse{}, nil
}
