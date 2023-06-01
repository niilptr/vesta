package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func GetModuleAddress() string {
	return sdk.AccAddress(crypto.AddressHash([]byte(ModuleName))).String()
}

func NewTwin(name string, hash string, creator string) Twin {
	return Twin{
		Name:       name,
		Hash:       hash,
		Creator:    creator,
		LastUpdate: creator,
	}
}

func NewEmptyTrainingState() TrainingState {
	return TrainingState{
		Value:                           false,
		TwinName:                        "",
		StartTime:                       time.Time{},
		TrainingConfigurationHash:       "",
		TrainingPhaseEndedConfirmations: make(map[string]bool),
		ValidationState: &ValidationState{
			Value:                       false,
			StartTime:                   time.Time{},
			MapValidatorsBestresulthash: make(map[string]string),
		},
	}
}
