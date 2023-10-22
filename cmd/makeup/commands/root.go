package commands

import (
	"github.com/cohix/makeup/pkg/makefile"
	"github.com/pkg/errors"
)

// Root is the root command
func Root(args []string) error {
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

	if err := mainmk.RunAll(); err != nil {
		return errors.Wrap(err, "failed to RunAll")
	}

	return nil
}
