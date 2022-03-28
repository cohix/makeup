package makefile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"

	"github.com/pkg/errors"

	"github.com/cohix/makeup/exec"
)

const (
	checkPrefix   = "# check "
	equalPrefix   = "# equal "
	includePrefix = "include "
)

// Makefile is a lightly-parsed Makefile
type Makefile struct {
	Checks   []Check
	Includes []string
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

func parse(file *os.File) (*Makefile, error) {
	mk := &Makefile{
		Checks:   []Check{},
		Includes: []string{},
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

			mk.Includes = append(mk.Includes, includePath)
		}
	}

	return mk, nil
}

type makeScanner struct {
	scn scanner.Scanner
}

func newScanner(rd io.Reader) *makeScanner {
	scn := scanner.Scanner{}
	scn.Init(rd)

	m := &makeScanner{
		scn: scn,
	}

	return m
}

// readLine reads the next line of the file
func (m *makeScanner) readLine() (string, error) {
	var err error

	m.scn.Error = func(_ *scanner.Scanner, msg string) {
		if msg != "" {
			err = errors.New(msg)
		}
	}

	buf := bytes.Buffer{}

	eof := false

	for {
		next := m.scn.Next()
		if next == scanner.EOF {
			eof = true
			break
		}

		stringNext := string(next)

		if err != nil {
			return "", errors.Wrap(err, "failed to scn.Next")
		}

		if stringNext == "\n" {
			break
		}

		if _, err := buf.Write([]byte(stringNext)); err != nil {
			return "", errors.Wrap(err, "failed to buf.Write")
		}
	}

	out := string(buf.Bytes())

	// skip empty lines, recursively
	if !eof && strings.TrimSpace(out) == "" {
		return m.readLine()
	}

	return out, nil
}
