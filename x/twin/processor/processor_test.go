package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	processortest "vesta/testutil/processor"
)

const testTwinName = "eva00"
const testTrainerMoniker = "val1"

func TestGetAccessToken(t *testing.T) {

	p := processortest.NewTestProcessor(t)

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	require.NotEmpty(t, acctoken)
}

func TestReadTrainConfiguration(t *testing.T) {

	p := processortest.NewTestProcessor(t)

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	tdc, twinRemoteURL, remoteURL, err := p.ReadTrainConfiguration(acctoken, testTwinName)
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
	isResultValid, reasonWhyFalse, err := p.ValidateTrainingResult(testTwinName, testTrainerMoniker)
	require.NoError(t, err)
	require.True(t, isResultValid)
	require.Empty(t, reasonWhyFalse)

}
