package types_test

import (
	"testing"
	"vesta/x/twin/types"

	"github.com/stretchr/testify/require"
)

func TestCheckMajorityAgreesOnTrainingBestResult(t *testing.T) {

	ts := types.NewEmptyTrainingState()
	ts.ValidationState.MapValidatorsBestresulthash["alice"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["bob"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["carol"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["dave"] = "abcde"

	agreement, twinHash := types.CheckMajorityAgreesOnTrainingBestResult(ts, 4)
	require.True(t, agreement)
	require.Equal(t, "abcde", twinHash)

	ts = types.NewEmptyTrainingState()
	ts.ValidationState.MapValidatorsBestresulthash["alice"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["bob"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["carol"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["dave"] = "12345"

	agreement, twinHash = types.CheckMajorityAgreesOnTrainingBestResult(ts, 4)
	require.True(t, agreement)
	require.Equal(t, "abcde", twinHash)

	ts = types.NewEmptyTrainingState()
	ts.ValidationState.MapValidatorsBestresulthash["alice"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["bob"] = "abcde"
	ts.ValidationState.MapValidatorsBestresulthash["carol"] = "12345"
	ts.ValidationState.MapValidatorsBestresulthash["dave"] = "12345"

	agreement, twinHash = types.CheckMajorityAgreesOnTrainingBestResult(ts, 4)
	require.False(t, agreement)
}
