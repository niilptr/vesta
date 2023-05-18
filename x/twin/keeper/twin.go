package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"vesta/x/twin/types"
)

// SetTwin set a specific twin in the store from its index
func (k Keeper) SetTwin(ctx sdk.Context, twin types.Twin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TwinKeyPrefix))
	b := k.cdc.MustMarshal(&twin)
	store.Set(types.TwinKey(
		twin.Name,
	), b)
}

// GetTwin returns a twin from its index
func (k Keeper) GetTwin(
	ctx sdk.Context,
	name string,

) (val types.Twin, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TwinKeyPrefix))

	b := store.Get(types.TwinKey(
		name,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTwin removes a twin from the store
func (k Keeper) RemoveTwin(
	ctx sdk.Context,
	name string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TwinKeyPrefix))
	store.Delete(types.TwinKey(
		name,
	))
}

// GetAllTwin returns all twin
func (k Keeper) GetAllTwin(ctx sdk.Context) (list []types.Twin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TwinKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Twin
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
