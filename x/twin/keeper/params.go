package keeper

import (
	"time"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAuthorizedAccounts(ctx sdk.Context) []string {

	var res []string
	k.paramstore.Get(ctx, types.KeyPrefix(types.KeyAuthorizedAccounts), &res)

	return res
}

func (k Keeper) GetMaxWaitingTraining(ctx sdk.Context) time.Duration {

	var res time.Duration
	k.paramstore.Get(ctx, types.KeyPrefix(types.KeyMaxWaitingTraining), &res)

	return res
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {

	return types.Params{
		AuthorizedAccounts: k.GetAuthorizedAccounts(ctx),
		MaxWaitingTraining: k.GetMaxWaitingTraining(ctx),
	}
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
