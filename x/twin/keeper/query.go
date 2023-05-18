package keeper

import (
	"vesta/x/twin/types"
)

var _ types.QueryServer = Keeper{}
