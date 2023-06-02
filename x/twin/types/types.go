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

func CheckMajorityAgreesOnTrainingPhaseEnded(ts TrainingState, maxConfirmations uint32) bool {

	count := 0
	for _, value := range ts.TrainingPhaseEndedConfirmations {
		if value == true {
			count++
		}
	}

	if float32(count) < float32(maxConfirmations*2/3) {
		return false
	}

	return true
}

func CheckMajorityAgreesOnTrainingBestResult(ts TrainingState, maxConfirmations uint32) (agreement bool, twinHash string) {

	countMap := make(map[string]uint32)

	for _, hash := range ts.ValidationState.MapValidatorsBestresulthash {
		countMap[hash] = countMap[hash] + 1
	}

	var maxCount uint32 = 0
	mostReputableHash := ""

	for hash, count := range countMap {
		if count > maxCount {
			maxCount = count
			mostReputableHash = hash
		}
	}

	if float32(maxCount) < float32(maxConfirmations)*2/3 {
		return false, mostReputableHash
	}

	return true, mostReputableHash
}
