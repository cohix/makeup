package makefile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/cohix/makeup/exec"
	"github.com/pkg/errors"
)

// RunAll runs all of the project components
func (m *Makefile) RunAll() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	binBase := filepath.Join(cwd, ".bin")

	errGroup, _ := errgroup.WithContext(context.Background())

	for _, incl := range m.Includes {
		componentDir := filepath.Dir(incl.Path)
		componentMakefile := filepath.Base(incl.Path)
		componentName := strings.TrimSuffix(componentMakefile, ".mk")

		fmt.Println("running:", componentName)

		binDest := filepath.Join(binBase, componentName)

		componentEnv, err := m.envForMkPath(incl.Path)
		if err != nil {
			return errors.Wrapf(err, "failed to envForMkPath %s", incl.Path)
		}

		// grab the 'env' target output and add some makeup-specific things
		env := append(
			strings.Split(componentEnv, "\n"),
			[]string{
				fmt.Sprintf("BIN_DEST=%s", binDest),
			}...,
		)

		errGroup.Go(func() error {
			writer := exec.NewPrefixWriter(componentName, os.Stdout)

			if _, err := exec.RunInDir(fmt.Sprintf("make -s -f %s run", componentMakefile), componentDir, writer, env...); err != nil {
				return errors.Wrapf(err, "failed to run %s", componentDir)
			}

			return nil
		})
	}

	return errGroup.Wait()
}

func (m *Makefile) envForMkPath(mkPath string) (string, error) {
	componentDir := filepath.Dir(mkPath)
	componentMakefile := filepath.Base(mkPath)
	componentName := strings.TrimSuffix(componentMakefile, ".mk")

	var out string
	var err error

	if m.ContainsOverride(componentName, "env") {
		overrideTarget := fmt.Sprintf("%s/%s", componentName, "env")

		out, err = exec.RunSilent(fmt.Sprintf("make -s -f %s %s", m.FullPath, overrideTarget), "")
		if err != nil {
			return "", errors.Wrapf(err, "failed to get override env %s", componentDir)
		}
	} else {
		out, err = exec.RunSilent(fmt.Sprintf("make -s -f %s env", componentMakefile), componentDir)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get env %s", componentDir)
		}
	}

	envLines := []string{}

	// remove any lines that don't look like an env statement, i.e. KEY=VALUE
	outLines := strings.Split(out, "\n")
	for _, l := range outLines {
		if strings.Contains(l, "=") {
			envLines = append(envLines, l)
		}
	}

	return strings.Join(envLines, "\n"), nil
}
