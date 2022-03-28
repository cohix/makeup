package exec

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

// PrefixWriter writes each line written into it to `out` prefixed with `prefix | `
type PrefixWriter struct {
	prefix string
	out    io.Writer

	lock sync.Mutex
	buf  []byte
}

// NewPrefixWriter creates a new PrefixWriter
func NewPrefixWriter(prefix string, out io.Writer) *PrefixWriter {
	p := &PrefixWriter{
		prefix: prefix,
		out:    out,
		lock:   sync.Mutex{},
	}

	return p
}

// Write takes input bytes and seperates it into lines, writing each to `out`
func (p *PrefixWriter) Write(in []byte) (int, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	inCopy := make([]byte, len(in))
	copy(inCopy, in)

	if len(p.buf) == 0 {
		p.buf = inCopy
	} else {
		p.buf = append(p.buf, inCopy...)
	}

	if !strings.Contains(string(p.buf), "\n") {
		fmt.Println(string(p.buf))
		return len(in), nil
	} else if len(p.buf) == 0 {
		return len(in), nil
	}

	fullLine := false
	if p.buf[len(p.buf)-1] == []byte("\n")[0] {
		fullLine = true
	}

	lines := bytes.Split(p.buf, []byte("\n"))

	if fullLine {
		p.buf = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	}

	for _, l := range lines {
		spaces := strings.Repeat(" ", 12-len(p.prefix))
		prefixVal := fmt.Sprintf("%s%s| ", p.prefix, spaces)

		prefixedLine := append([]byte(prefixVal), l...)
		prefixedLine = append(prefixedLine, []byte("\n")...)

		p.out.Write(prefixedLine)
	}

	return len(in), nil
}
