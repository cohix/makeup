package commands

import (
	"github.com/cohix/makeup/pkg/makefile"
	"github.com/pkg/errors"
)

// Test runs a test on every component of the project
func Test(args []string) error {
	mainmk, err := makefile.Parse("./main.mk")
	if err != nil {
		return errors.Wrap(err, "failed to Parse main.mk")
	}

	if err := mainmk.TestChecks(); err != nil {
		return errors.Wrap(err, "failed to TestChecks")
	}

	if err := mainmk.TestAll(); err != nil {
		return errors.Wrap(err, "failed to TestAll")
	}

	return nil
}
