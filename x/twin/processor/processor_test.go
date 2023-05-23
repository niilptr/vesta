package processor_test

import (
	"fmt"
	"os"
	"strings"
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
	lastok := strings.Compare(userHome[len(userHome)-1:], "/")
	if lastok != 0 {
		userHome = userHome + "/"
	}
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

	acctoken, err := p.GetAccessToken()
	require.NoError(t, err)
	content, err := p.ReadTrainConfiguration(acctoken)
	require.NoError(t, err)
	require.NotEmpty(t, content)
}
