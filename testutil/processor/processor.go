package processor

import (
	"fmt"
	"os"
	"testing"
	keepertest "vesta/testutil/keeper"
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	"github.com/stretchr/testify/require"
)

const PathFromHomeToTestDir = "test-vesta/"

func NewTestProcessor(t *testing.T) processor.Processor {

	userHome, err := os.UserHomeDir()
	require.NoError(t, err)
	userHome = processor.CheckPathFormat(userHome)
	_, ctx := keepertest.NewTestKeeper(t)
	l := ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))

	return processor.NewProcessor(userHome+PathFromHomeToTestDir, l)
}
