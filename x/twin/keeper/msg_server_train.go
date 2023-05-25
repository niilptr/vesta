package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ms msgServer) Train(goCtx context.Context, msg *types.MsgTrain) (*types.MsgTrainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := ms.Keeper.StartTraining(ctx, msg.Name, msg.Creator, msg.TrainingConfigurationHash)
	if err != nil {
		return &types.MsgTrainResponse{}, err
	}

	return &types.MsgTrainResponse{}, nil
}
