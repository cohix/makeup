package commands

import (
	"github.com/cohix/makeup/pkg/makefile"
	"github.com/pkg/errors"
)

// Build builds every component of the project
func Build(args []string) error {
	mainmk, err := makefile.Parse("./main.mk")
	if err != nil {
		return errors.Wrap(err, "failed to Parse main.mk")
	}

	if err := mainmk.TestChecks(); err != nil {
		return errors.Wrap(err, "failed to TestChecks")
	}

	if err := mainmk.BuildAll(); err != nil {
		return errors.Wrap(err, "failed to BuildAll")
	}

	return nil
}
