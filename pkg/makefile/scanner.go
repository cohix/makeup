package makefile

import (
	"bytes"
	"io"
	"strings"
	"text/scanner"

	"github.com/pkg/errors"
)

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
