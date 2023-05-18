package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (ms msgServer) CreateTwin(goCtx context.Context, msg *types.MsgCreateTwin) (*types.MsgCreateTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, found := ms.Keeper.GetTwin(ctx, msg.Name)
	if found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	twin := types.NewTwin(msg.Name, msg.Hash, msg.Creator)

	ms.Keeper.SetTwin(ctx, twin)
	return &types.MsgCreateTwinResponse{}, nil
}

func (ms msgServer) UpdateTwin(goCtx context.Context, msg *types.MsgUpdateTwin) (*types.MsgUpdateTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	twin, found := ms.Keeper.GetTwin(ctx, msg.Name)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != twin.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	twin = types.NewTwin(msg.Name, msg.Hash, msg.Creator)

	ms.Keeper.SetTwin(ctx, twin)

	return &types.MsgUpdateTwinResponse{}, nil
}

func (ms msgServer) DeleteTwin(goCtx context.Context, msg *types.MsgDeleteTwin) (*types.MsgDeleteTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	twin, found := ms.Keeper.GetTwin(ctx, msg.Name)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != twin.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	ms.Keeper.RemoveTwin(ctx, msg.Name)

	return &types.MsgDeleteTwinResponse{}, nil
}
