package makefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/cohix/makeup/exec"
)

const (
	includePrefix = "include "
	checkPrefix   = "# check "
	equalPrefix   = "# equal "
	externPrefix  = "# extern "
	rootPrefix    = "# root "
	overrideLine  = "# override"
)

// Makefile is a lightly-parsed Makefile
type Makefile struct {
	Checks    []Check
	Includes  []include
	Overrides []override

	FullPath string
}

// include represents an `include` statement in a Makefile, plus optional `extern` modifier
type include struct {
	Path   string
	Extern string
}

// override represents an overridden target for a component
type override struct {
	Component string
	Target    string
}

// Parse reads and parses the Makefile at the given path
func Parse(path string) (*Makefile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to Open %s", path)
	}

	defer file.Close()

	mk, err := parse(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", path)
	}

	if err := mk.ensureIncludes(); err != nil {
		return nil, errors.Wrap(err, "failed to ensureIncludes")
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to filepath.Abs")
	}

	mk.FullPath = fullPath

	return mk, nil
}

// TestChecks runs each defined Check and returns an error if any fail
func (m *Makefile) TestChecks() error {
	for _, c := range m.Checks {
		out, err := exec.RunSilent(c.Cmd, "")
		if err != nil {
			return errors.Wrapf(err, "failed to RunSilent %s", c.Cmd)
		}

		if !strings.Contains(out, c.Equals) {
			return fmt.Errorf("failed check: %s is not %s, got %s", c.Cmd, c.Equals, out)
		}
	}

	return nil
}

// ContainsOverride returns true if the main.mk contains an overridden target for the given component
func (m *Makefile) ContainsOverride(component, target string) bool {
	for _, o := range m.Overrides {
		if o.Component == component && o.Target == target {
			return true
		}
	}

	return false
}

func parse(file *os.File) (*Makefile, error) {
	mk := &Makefile{
		Checks:   []Check{},
		Includes: []include{},
	}

	scn := newScanner(file)

	for {
		line, err := scn.readLine()
		if err != nil {
			return nil, errors.Wrap(err, "failed to readLine")
		}

		if len(line) == 0 {
			break
		}

		if strings.HasPrefix(line, checkPrefix) {
			check := Check{
				Cmd: strings.TrimPrefix(line, checkPrefix),
			}

			nextLine, err := scn.readLine()
			if err != nil {
				return nil, errors.Wrap(err, "failed to readLine")
			}

			if !strings.HasPrefix(nextLine, equalPrefix) {
				return nil, fmt.Errorf("line following check is not an 'equal' value (got %s)", nextLine)
			}

			check.Equals = strings.TrimPrefix(nextLine, equalPrefix)

			mk.Checks = append(mk.Checks, check)
		} else if strings.HasPrefix(line, includePrefix) {
			includePath := strings.TrimPrefix(line, includePrefix)

			incl := include{
				Path: includePath,
			}

			mk.Includes = append(mk.Includes, incl)
		} else if strings.HasPrefix(line, externPrefix) {
			externPath := strings.TrimPrefix(line, externPrefix)

			includeLine, err := scn.readLine()
			if err != nil {
				return nil, errors.Wrap(err, "failed to readLine")
			}

			if !strings.HasPrefix(includeLine, includePrefix) {
				return nil, fmt.Errorf("line following extern is not an 'include' statement (got %s)", includeLine)
			}

			includePath := strings.TrimPrefix(includeLine, includePrefix)

			incl := include{
				Path:   includePath,
				Extern: externPath,
			}

			mk.Includes = append(mk.Includes, incl)
		} else if line == overrideLine {
			targetLine, err := scn.readLine()
			if err != nil {
				return nil, errors.Wrap(err, "failed to readLine")
			}

			if !strings.Contains(targetLine, ":") {
				return nil, fmt.Errorf("line following override is not a target (got %s)", targetLine)
			}

			fullTarget := targetLine[:strings.Index(targetLine, ":")]
			targetParts := strings.Split(fullTarget, "/")
			if len(targetParts) != 2 {
				return nil, fmt.Errorf("override targed must have two /-seperated parts (got %d)", len(targetParts))
			}

			component := targetParts[0]
			target := targetParts[1]

			ovr := override{
				Component: component,
				Target:    target,
			}

			mk.Overrides = append(mk.Overrides, ovr)
		}
	}

	return mk, nil
}

func (m *Makefile) ensureIncludes() error {
	for _, incl := range m.Includes {
		if _, err := os.Stat(incl.Path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if incl.Extern != "" {
					return errors.Wrapf(err, "missing %s from extern %s", incl.Path, incl.Extern)
				} else {
					return errors.Wrap(err, "missing %s")
				}
			}

			return errors.Wrapf(err, "failed to Stat %s", incl.Path)
		}
	}

	return nil
}
