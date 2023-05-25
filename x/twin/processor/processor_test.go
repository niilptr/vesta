package processor_test

import (
	"fmt"
	"os"
	"testing"
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	keepertest "vesta/testutil/keeper"

	"github.com/stretchr/testify/require"
)

const PathFromHomeToTestDir = "test-vesta/"

func NewTestProcessor(t *testing.T) processor.Processor {

	userHome, err := os.UserHomeDir()
	require.NoError(t, err)
	userHome = processor.CheckPathFormat(userHome)
	_, ctx := keepertest.TwinKeeper(t)
	l := ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))

	return processor.NewProcessor(userHome+PathFromHomeToTestDir, l)
}

func TestGetAccessToken(t *testing.T) {

	p := NewTestProcessor(t)

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	require.NotEmpty(t, acctoken)
}

func TestReadTrainConfiguration(t *testing.T) {

	p := NewTestProcessor(t)
	twinName := "eva00"

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	tdc, twinRemoteURL, remoteURL, err := p.ReadTrainConfiguration(acctoken, twinName)
	require.NoError(t, err)
	require.NotEmpty(t, tdc)
	require.NotEmpty(t, twinRemoteURL)
	require.NotEmpty(t, remoteURL)
}
