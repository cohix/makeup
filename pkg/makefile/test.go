package makefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cohix/makeup/pkg/exec"
	"github.com/pkg/errors"
)

// TestAll sequentially runs each of the project components' test targets
func (m *Makefile) TestAll() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	binBase := filepath.Join(cwd, ".bin")

	for _, incl := range m.Includes {
		componentDir := filepath.Dir(incl.Path)
		componentMakefile := filepath.Base(incl.Path)
		componentName := strings.TrimSuffix(componentMakefile, ".mk")

		fmt.Println("testing:", componentName)

		binDest := filepath.Join(binBase, componentName)

		env := []string{
			fmt.Sprintf("BIN_DEST=%s", binDest),
		}

		if _, err := exec.RunInDir(fmt.Sprintf("make -s -f %s test", componentMakefile), componentDir, nil, env...); err != nil {
			return errors.Wrapf(err, "failed to test %s", componentDir)
		}

		fmt.Println("test complete:", componentName)
	}

	return nil
}
