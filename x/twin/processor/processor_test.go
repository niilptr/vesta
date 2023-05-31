package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	processortest "vesta/testutil/processor"
)

const testTwinNameTraining = "eva00"
const testTwinNameValidation = "eva00"
const testTrainerMoniker = "val1"

func TestNewTestProcessor(t *testing.T) {
	_, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)
}

func TestReadTrainConfiguration(t *testing.T) {

	p, err := processortest.NewTestProcessor(t)
	require.NoError(t, err)

	tdc, twinRemoteURL, err := p.ReadTrainConfiguration(p.GetAccessToken(), testTwinNameTraining)
	require.NoError(t, err)
	require.NotEmpty(t, tdc)
	require.NotEmpty(t, twinRemoteURL)
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

	vtr, err := p.ReadValidatorsTrainingResults(testTwinNameValidation)
	_, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
	isResultValid, reasonWhyFalse, err := p.ValidateBestTrainingResult(testTwinNameValidation, trainerMoniker, newTwinHash)
	require.NoError(t, err)
	require.True(t, isResultValid)
	require.Empty(t, reasonWhyFalse)

}
