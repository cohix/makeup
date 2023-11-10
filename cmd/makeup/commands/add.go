package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const tmpl = `
build:
	go build -o ${BIN_DEST}

run:
	${BIN_DEST}

test:
	go test -v ./...

env:
	echo "CONFIG_KEY=some_config_val"

clean:
	rm ${BIN_DEST}
`

func Add(args []string) error {
	if len(args) < 1 {
		return errors.New("missing arg: component name")
	}

	componentArg := strings.ToLower(args[0])
	componentArgParts := strings.Split(componentArg, fmt.Sprintf("%c", filepath.Separator))
	componentName := componentArgParts[len(componentArgParts)-1]

	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to Getwd")
	}

	componentDir := filepath.Join(wd, componentArg)

	_, err = os.Stat(componentDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if mkErr := os.MkdirAll(componentDir, os.ModePerm); mkErr != nil {
				return errors.Wrapf(err, "failed to create component dir %s", componentDir)
			}
		} else {
			return errors.Wrapf(err, "failed to os.Stat %s", componentDir)
		}
	}

	filename := fmt.Sprintf("%s.mk", componentName)
	componentMkFilepath := filepath.Join(componentDir, filename)

	if err := os.WriteFile(componentMkFilepath, []byte(tmpl), os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to os.WriteFile %s", componentMkFilepath)
	}

	relativeComponentMkFilepath, err := filepath.Rel(wd, componentMkFilepath)
	if err != nil {
		return errors.Wrap(err, "failed to filepath.Rel")
	}

	includeStatement := fmt.Sprintf("\ninclude ./%s", relativeComponentMkFilepath)

	mainMkFilepath := filepath.Join(wd, "main.mk")

	_, err = os.Stat(mainMkFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(mainMkFilepath, []byte(includeStatement), os.ModePerm); err != nil {
				return errors.Wrapf(err, "failed to os.WriteFile (new) %s", mainMkFilepath)
			}
		} else {
			return errors.Wrapf(err, "failed to os.Stat %s", mainMkFilepath)
		}
	} else {
		mkFile, err := os.OpenFile(mainMkFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "failed to os.OpenFile %s", mainMkFilepath)
		}

		defer mkFile.Close()

		if _, err = mkFile.WriteString(includeStatement); err != nil {
			return errors.Wrapf(err, "failed to WriteString %s", mainMkFilepath)
		}
	}

	fmt.Println("component created:", componentMkFilepath)

	return nil
}
