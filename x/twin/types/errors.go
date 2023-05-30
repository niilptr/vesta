package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/twin module sentinel errors
var (

	// Chain errors
	ErrTrainingStateNotFound           = sdkerrors.Register(ModuleName, 1101, "training state informations not found")
	ErrTrainingInProgress              = sdkerrors.Register(ModuleName, 1102, "a training is in progress")
	ErrTrainingNotInProgress           = sdkerrors.Register(ModuleName, 1103, "no training is in progress")
	ErrTrainingValidationNotInProgress = sdkerrors.Register(ModuleName, 1104, "no training validation is in progress")
)
