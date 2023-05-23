package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TrainingState(goCtx context.Context, req *types.QueryGetTrainingStateRequest) (*types.QueryGetTrainingStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetTrainingState(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTrainingStateResponse{TrainingState: val}, nil
}
