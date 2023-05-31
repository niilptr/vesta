package processor

import (
	"fmt"
	"os"
	"testing"
	keepertest "vesta/testutil/keeper"
	"vesta/x/twin/processor"
	"vesta/x/twin/types"
)

const PathFromHomeToTestDir = "test-vesta/"

func NewTestProcessor(t *testing.T) (p processor.Processor, err error) {

	userHome, err := os.UserHomeDir()
	if err != nil {
		return p, err
	}

	userHome = processor.CheckPathFormat(userHome)
	_, ctx := keepertest.NewTestKeeper(t)

	l := ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))

	p, err = processor.NewProcessor(userHome+PathFromHomeToTestDir, l)
	if err != nil {
		return p, err
	}

	return p, nil
}
