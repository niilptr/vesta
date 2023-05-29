package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	processortest "vesta/testutil/processor"
)

func TestGetAccessToken(t *testing.T) {

	p := processortest.NewTestProcessor(t)

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	require.NotEmpty(t, acctoken)
}

func TestReadTrainConfiguration(t *testing.T) {

	p := processortest.NewTestProcessor(t)
	twinName := "eva00"

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	tdc, twinRemoteURL, remoteURL, err := p.ReadTrainConfiguration(acctoken, twinName)
	require.NoError(t, err)
	require.NotEmpty(t, tdc)
	require.NotEmpty(t, twinRemoteURL)
	require.NotEmpty(t, remoteURL)
}
