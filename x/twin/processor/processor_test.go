package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	processortest "vesta/testutil/processor"
)

const testTwinName = "eva00"
const testTrainerMoniker = "val1"

func TestNewTestProcessor(t *testing.T) {
	_ = processortest.NewTestProcessor(t)
}

func TestReadTrainConfiguration(t *testing.T) {

	p := processortest.NewTestProcessor(t)

	tdc, twinRemoteURL, remoteURL, err := p.ReadTrainConfiguration(p.GetAccessToken(), testTwinName)
	require.NoError(t, err)
	require.NotEmpty(t, tdc)
	require.NotEmpty(t, twinRemoteURL)
	require.NotEmpty(t, remoteURL)
}

func TestTrain(t *testing.T) {
	p := processortest.NewTestProcessor(t)
	err := p.Train()
	require.NoError(t, err)
}

func TestValidateTrainingResult(t *testing.T) {
	p := processortest.NewTestProcessor(t)
	vtr, err := p.ReadValidatorsTrainingResults(testTwinName)
	_, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
	isResultValid, reasonWhyFalse, err := p.ValidateBestTrainingResult(testTwinName, trainerMoniker, newTwinHash)
	require.NoError(t, err)
	require.True(t, isResultValid)
	require.Empty(t, reasonWhyFalse)

}
