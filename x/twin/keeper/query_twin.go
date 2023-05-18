package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vesta/x/twin/types"
)

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
