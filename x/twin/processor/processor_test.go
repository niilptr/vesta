package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	processortest "vesta/testutil/processor"
)

const (
	testTwinNameTraining   = "eva00"
	testTwinNameValidation = "eva00"
	testTrainerMoniker     = "val1"
	testTrainConfHash      = "a447bde2cd1a95e421fb33567be40c89940e2b4d572f7ca0f73826b5838108d8"
)

func TestNewTestProcessor(t *testing.T) {
	_, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)
}

func TestReadTrainConfigurationAndVerifyHash(t *testing.T) {

	p, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)

	tdc, twinRemoteURL, err := p.ReadTrainConfigurationAndVerifyHash(testTwinNameTraining, testTrainConfHash)
	require.NoError(t, err)
	require.NotEmpty(t, tdc)
	require.NotEmpty(t, twinRemoteURL)

	tdc, twinRemoteURL, err = p.ReadTrainConfigurationAndVerifyHash(testTwinNameTraining, "not a valid hash")
	require.Error(t, err)
	require.Empty(t, tdc)
	require.Empty(t, twinRemoteURL)
}

func TestTrain(t *testing.T) {
	p, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)

	err = p.Train()
	require.NoError(t, err)
}

func TestValidateTrainingResult(t *testing.T) {
	p, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)

	vtr, err := p.ReadValidatorsTrainingResults(testTwinNameValidation, "not a valid hash")
	require.Error(t, err)

	vtr, err = p.ReadValidatorsTrainingResults(testTwinNameValidation, testTrainConfHash)
	require.NoError(t, err)

	_, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
	isResultValid, reasonWhyFalse, err := p.ValidateBestTrainingResult(testTwinNameValidation, trainerMoniker, newTwinHash)
	require.NoError(t, err)
	require.True(t, isResultValid)
	require.Empty(t, reasonWhyFalse)

}
