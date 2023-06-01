package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"vesta/x/twin/processor"
	"vesta/x/twin/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		nodeHome   string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	nodeHome string,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		nodeHome:   nodeHome,
	}
}

func (k Keeper) GetNodeHome() string {
	return k.nodeHome
}

func ModuleLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ModuleLogger(ctx)
}

// ====================================================================================
// Params
// ====================================================================================

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

// ====================================================================================
// Twin
// ====================================================================================

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

func (k Keeper) UpdateTwinFromVestaTraining(ctx sdk.Context, name string, hash string) {

	twin, found := k.GetTwin(ctx, name)
	if !found {
		k.SetTwin(ctx, types.NewTwin(name, hash, types.ModuleName))
	}

	twin.LastUpdate = types.GetModuleAddress()
	k.SetTwin(ctx, twin)
}

// ====================================================================================
// Training state
// ====================================================================================

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) SetTrainingState(ctx sdk.Context, trainingState types.TrainingState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	b := k.cdc.MustMarshal(&trainingState)
	store.Set([]byte{0}, b)
}

// UpdateTrainingStateValue set trainingState value in the store
func (k Keeper) MustUpdateTrainingStateValue(ctx sdk.Context, ts types.TrainingState, newValue bool) types.TrainingState {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.Value = newValue
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
	ts, found := k.GetTrainingState(ctx)
	if !found {
		panic("Training state not found after its updating.")
	}
	return ts
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) MustUpdateTrainingStateValidationValue(ctx sdk.Context, ts types.TrainingState, newValue bool) types.TrainingState {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.ValidationState.Value = newValue
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
	ts, found := k.GetTrainingState(ctx)
	if !found {
		panic("Training state not found after its updating.")
	}
	return ts
}

// GetTrainingState returns trainingState
func (k Keeper) GetTrainingState(ctx sdk.Context) (ts types.TrainingState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return ts, false
	}

	k.cdc.MustUnmarshal(b, &ts)
	return ts, true
}

// GetTraining returns trainingState
func (k Keeper) GetTrainingStateValue(ctx sdk.Context) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return false
	}
	ts := types.TrainingState{}
	k.cdc.MustUnmarshal(b, &ts)

	return ts.Value
}

// GetTraining returns trainingState
func (k Keeper) GetTrainingStateTwinName(ctx sdk.Context) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return ""
	}
	ts := types.TrainingState{}
	k.cdc.MustUnmarshal(b, &ts)

	return ts.TwinName
}

// RemoveTrainingState removes trainingState from the store
func (k Keeper) RemoveTrainingState(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	store.Delete([]byte{0})
}

// ====================================================================================
// Train
// ====================================================================================

func (k Keeper) StartTraining(ctx sdk.Context, twinName string, creator string, trainConfHash string) error {

	isTraining := k.GetTrainingStateValue(ctx)

	if isTraining {
		return types.ErrTrainingInProgress
	}

	// Keeper acts before processor beacuse processor methods can lead to
	// non-deterministic results (due e.g. to problems reaching the
	// central db)
	k.SetTrainingState(ctx, types.TrainingState{
		Value:                     true,
		TwinName:                  twinName,
		StartTime:                 ctx.BlockTime(),
		TrainingConfigurationHash: trainConfHash,
	})

	////////// START GO ROUTINE ////////////////
	// Processor will:
	// 1. get the training configuration from remote (that configuration contains all
	//    trainers specific configuration);
	// 2. verify training configuration match the one provided;
	// 3. select the specific train configuration to run on the Vesta node;
	// 4. run the local training process
	//
	// The local training will:
	// 1. get the specific training configuration;
	// 2. train the twin model;
	// 3. upload the training results.
	go processor.StartProcessorForTrainingTwin(k.GetNodeHome(), k.Logger(ctx), twinName, trainConfHash)

	return nil
}

// ====================================================================================
// Confirm train phase ended
// ====================================================================================

// Each authorized account has to check if training phase ended and broadcast its
// confirmation. This confirmation will be stored in the training state, so later
// it will be possible to verify if majority agrees on this.
func (k Keeper) AddTrainingPhaseEndedConfirmation(ctx sdk.Context, signer string) error {

	ts, found := k.GetTrainingState(ctx)

	// Cannot be possible that training state is not found because a train request would have
	// initialized it.
	if !found {
		return types.ErrTrainingStateNotFound
	}

	// Cannot be possible that training state value (aka isTraining) is not set to true
	// (because it will be modified after majority of confirmations is enstablished).
	if !ts.Value {
		return types.ErrTrainingNotInProgress
	}

	ts.TrainingPhaseEndedConfirmations[signer] = true

	k.SetTrainingState(ctx, ts)

	return nil
}

func (k Keeper) CheckMajorityAgreesOnTrainingPhaseEnded(ctx sdk.Context, ts types.TrainingState, maxConfirmations uint32) bool {

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

// ====================================================================================
// Confirm best result is
// ====================================================================================

// Each authorized account has to broadcast its best result, and this will be stored in
// the training state, so later it will be possible to verify if majority agrees on a
// specific result. Best result is actually represented by its hash (the new twin hash).
func (k Keeper) AddBestTrainResultToTrainingState(ctx sdk.Context, signer string, twinHash string) error {

	ts, found := k.GetTrainingState(ctx)

	// Cannot be possible that training state is not found because a train request
	// would have initialized it.
	if !found {
		return types.ErrTrainingStateNotFound
	}

	// Cannot be possible that training state value (aka isTraining) is set to true
	// (because validation phase comes after training phase ended).
	if ts.Value {
		return types.ErrTrainingInProgress
	}

	// Cannot be possible that validation state value (aka isValidating) is set to
	// false (because validation phase must be active before agreement on best result
	// is reached).
	if !ts.ValidationState.Value {
		return types.ErrTrainingValidationNotInProgress
	}

	ts.ValidationState.MapValidatorsBestresulthash[signer] = twinHash

	k.SetTrainingState(ctx, ts)

	return nil
}

func (k Keeper) CheckMajorityAgreesOnTrainingBestResult(ctx sdk.Context, ts types.TrainingState, maxConfirmations uint32) (agreement bool, twinHash string) {

	countMap := make(map[string]uint32)

	for key := range ts.ValidationState.MapValidatorsBestresulthash {
		countMap[key] = countMap[key] + 1
	}

	var maxCount uint32 = 0
	mostReputableHash := ""

	for hash, count := range countMap {
		if count > maxCount {
			maxCount = count
			mostReputableHash = hash
		}
	}

	if float32(maxCount) < float32(maxConfirmations*2/3) {
		return false, mostReputableHash
	}

	return true, mostReputableHash
}

// ====================================================================================
// Authorization
// ====================================================================================

func (k Keeper) IsAccountAuthorized(ctx sdk.Context, address string) (bool, error) {

	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	for _, addr := range k.GetAuthorizedAccounts(ctx) {
		if address == addr {
			return true, nil
		}
	}

	return false, nil

}
