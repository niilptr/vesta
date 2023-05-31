package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAuthorizedAccounts(ctx sdk.Context) []string {
	return k.GetParams(ctx).AuthorizedAccounts
}
