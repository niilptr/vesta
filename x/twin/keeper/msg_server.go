package keeper

import (
	"context"

	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// ====================================================================================
// Twin
// ====================================================================================

func (ms msgServer) CreateTwin(goCtx context.Context, msg *types.MsgCreateTwin) (*types.MsgCreateTwinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgCreateTwinResponse{}, err
	}

	if !authorized {
		return &types.MsgCreateTwinResponse{}, types.ErrAccountNotAuthorized
	}

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

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgUpdateTwinResponse{}, err
	}

	if !authorized {
		return &types.MsgUpdateTwinResponse{}, types.ErrAccountNotAuthorized
	}

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

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgDeleteTwinResponse{}, err
	}

	if !authorized {
		return &types.MsgDeleteTwinResponse{}, types.ErrAccountNotAuthorized
	}

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

// ====================================================================================
// Train
// ====================================================================================

func (ms msgServer) Train(goCtx context.Context, msg *types.MsgTrain) (*types.MsgTrainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgTrainResponse{}, err
	}

	if !authorized {
		return &types.MsgTrainResponse{}, types.ErrAccountNotAuthorized
	}

	err = ms.Keeper.StartTraining(ctx, msg.Name, msg.Creator, msg.TrainingConfigurationHash)
	if err != nil {
		return &types.MsgTrainResponse{}, err
	}

	return &types.MsgTrainResponse{}, nil
}

// ====================================================================================
// Confirm train phase ended
// ====================================================================================

func (ms msgServer) ConfirmTrainPhaseEnded(goCtx context.Context, msg *types.MsgConfirmTrainPhaseEnded) (*types.MsgConfirmTrainPhaseEndedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgConfirmTrainPhaseEndedResponse{}, err
	}

	if !authorized {
		return &types.MsgConfirmTrainPhaseEndedResponse{}, types.ErrAccountNotAuthorized
	}

	err = ms.Keeper.AddTrainingPhaseEndedConfirmation(ctx, msg.Creator)
	if err != nil {
		return &types.MsgConfirmTrainPhaseEndedResponse{}, err
	}

	return &types.MsgConfirmTrainPhaseEndedResponse{}, nil
}

// ====================================================================================
// Confirm best result is
// ====================================================================================

func (ms msgServer) ConfirmBestTrainResultIs(goCtx context.Context, msg *types.MsgConfirmBestTrainResultIs) (*types.MsgConfirmBestTrainResultIsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authorization
	authorized, err := ms.Keeper.IsAccountAuthorized(ctx, msg.Creator)
	if err != nil {
		return &types.MsgConfirmBestTrainResultIsResponse{}, err
	}

	if !authorized {
		return &types.MsgConfirmBestTrainResultIsResponse{}, types.ErrAccountNotAuthorized
	}

	err = ms.Keeper.AddBestTrainResultToTrainingState(ctx, msg.Creator, msg.Hash)
	if err != nil {
		return &types.MsgConfirmBestTrainResultIsResponse{}, err
	}

	return &types.MsgConfirmBestTrainResultIsResponse{}, nil
}
