package makefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cohix/makeup/exec"
	"github.com/pkg/errors"
)

// BuildAll sequentially runs each of the project components' build targets
func (m *Makefile) BuildAll() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	binBase := filepath.Join(cwd, ".bin")

	for _, mkPath := range m.Includes {
		componentDir := filepath.Dir(mkPath)
		componentMakefile := filepath.Base(mkPath)
		componentName := strings.TrimSuffix(componentMakefile, ".mk")

		fmt.Println("building:", componentName)

		binDest := filepath.Join(binBase, componentName)

		env := []string{
			fmt.Sprintf("BIN_DEST=%s", binDest),
		}

		if _, err := exec.RunInDir(fmt.Sprintf("make -s -f %s build", componentMakefile), componentDir, nil, env...); err != nil {
			return errors.Wrapf(err, "failed to build %s", componentDir)
		}

		fmt.Println("build complete:", componentName)
	}

	return nil
}
