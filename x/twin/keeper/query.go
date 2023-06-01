package keeper

import (
	"context"

	"vesta/x/twin/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// =====================================================================
// Params
// =====================================================================
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// =====================================================================
// Twin
// =====================================================================
func (k Keeper) TwinAll(goCtx context.Context, req *types.QueryAllTwinRequest) (*types.QueryAllTwinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var twins []types.Twin
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	twinStore := prefix.NewStore(store, types.KeyPrefix(types.TwinKeyPrefix))

	pageRes, err := query.Paginate(twinStore, req.Pagination, func(key []byte, value []byte) error {
		var twin types.Twin
		if err := k.cdc.Unmarshal(value, &twin); err != nil {
			return err
		}

		twins = append(twins, twin)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTwinResponse{Twin: twins, Pagination: pageRes}, nil
}

func (k Keeper) Twin(goCtx context.Context, req *types.QueryGetTwinRequest) (*types.QueryGetTwinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetTwin(
		ctx,
		req.Name,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTwinResponse{Twin: val}, nil
}

// =====================================================================
// Training state
// =====================================================================
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
