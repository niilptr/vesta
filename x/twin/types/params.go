package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	KeyAuthorizedAccounts = "AuthorizedAccounts"
	KeyMaxWaitingTraining = "MaxWaitingTraining"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		AuthorizedAccounts: []string{},
		MaxWaitingTraining: 60 * time.Second,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(KeyAuthorizedAccounts), &p.AuthorizedAccounts, validateAuthorizedAccounts),
		paramtypes.NewParamSetPair(KeyPrefix(KeyMaxWaitingTraining), &p.MaxWaitingTraining, validateMaxWaitingTraining),
	}
}

func validateAuthorizedAccounts(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, a := range v {
		_, err := sdk.AccAddressFromBech32(a)
		if err != nil {
			return fmt.Errorf("authorized address invalid Bech32: %s", a)
		}
	}

	return nil
}

func validateMaxWaitingTraining(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max waiting time must be positive: %d", v)
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	// TODO: define params validation
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
