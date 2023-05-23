package types

import (
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
