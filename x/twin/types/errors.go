package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/twin module sentinel errors
var (

	// Chain errors
	ErrTrainingInProgress = sdkerrors.Register(ModuleName, 1101, "a training is in progress")
)
