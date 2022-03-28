package commands

import (
	"github.com/cohix/makeup/makefile"
	"github.com/pkg/errors"
)

// Clean runs clean on every component of the project
func Clean(args []string) error {
	mainmk, err := makefile.Parse("./main.mk")
	if err != nil {
		return errors.Wrap(err, "failed to Parse main.mk")
	}

	if err := mainmk.TestChecks(); err != nil {
		return errors.Wrap(err, "failed to TestChecks")
	}

	if err := mainmk.CleanAll(); err != nil {
		return errors.Wrap(err, "failed to CleanAll")
	}

	return nil
}
