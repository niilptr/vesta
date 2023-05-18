package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"vesta/x/twin/types"
)

func (k msgServer) CreateTwin(goCtx context.Context, msg *types.MsgCreateTwin) (*types.MsgCreateTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetTwin(
		ctx,
		msg.Name,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var twin = types.Twin{
		Creator: msg.Creator,
		Name:    msg.Name,
		Hash:    msg.Hash,
	}

	k.SetTwin(
		ctx,
		twin,
	)
	return &types.MsgCreateTwinResponse{}, nil
}

func (k msgServer) UpdateTwin(goCtx context.Context, msg *types.MsgUpdateTwin) (*types.MsgUpdateTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetTwin(
		ctx,
		msg.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var twin = types.Twin{
		Creator: msg.Creator,
		Name:    msg.Name,
		Hash:    msg.Hash,
	}

	k.SetTwin(ctx, twin)

	return &types.MsgUpdateTwinResponse{}, nil
}

func (k msgServer) DeleteTwin(goCtx context.Context, msg *types.MsgDeleteTwin) (*types.MsgDeleteTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetTwin(
		ctx,
		msg.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveTwin(
		ctx,
		msg.Name,
	)

	return &types.MsgDeleteTwinResponse{}, nil
}
