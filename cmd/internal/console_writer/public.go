package console_writer

import (
	"io"
	"os"

	"github.com/joyent/triton-go/cmd/internal/pager"
	"github.com/rs/zerolog"
)

type ConsoleWriter interface {
	io.Writer
	io.Closer
	Wait() error
	getPager() *pager.Pager
}

var globalTerm ConsoleWriter
var origStdout *os.File
var origStderr *os.File

type terminal struct {
	stdout io.Writer
	stderr io.Writer
}

func GetTerminal() ConsoleWriter {
	return globalTerm
}

func UsePager(usePager bool) error {
	p := globalTerm.getPager()

	switch {
	case usePager && p != nil:
		return nil
	case usePager && p == nil:
		p, err := pager.New()
		if err != nil {
			return err
		}
		t := &terminal{
			stdout: p,
		}
		globalTerm = t
		return nil
	case !usePager && p != nil:
		t := &terminal{
			stdout: origStdout,
		}
		globalTerm = t
		return nil
	case !usePager && p == nil:
		t := &terminal{
			stdout: origStdout,
		}
		globalTerm = t
		return nil
	default:
		panic("x")
	}
	// if stdout is already a pager then nothing
	// if stdout is terminal then insert pager
	// if stdout is pager then replace with terminal
	// if stdout is a terminal then nothing
}

func (t *terminal) Close() error {
	if p := t.getPager(); p != nil {
		return p.Wait()
	}

	return nil
}

func (t *terminal) Wait() error {
	if p := t.getPager(); p != nil {
		return p.Wait()
	}

	return nil
}

func (t *terminal) getPager() *pager.Pager {
	if p, ok := t.stdout.(*pager.Pager); ok {
		return p
	}

	return nil
}
func (t *terminal) Write(data []byte) (n int, err error) {
	return t.stdout.Write(data)
}

func init() {
	// we need to keep copy of StdOut
	globalTerm = &terminal{
		stdout: zerolog.SyncWriter(os.Stdout),
		stderr: zerolog.SyncWriter(os.Stderr),
	}

	origStdout = os.Stdout
	origStderr = os.Stderr
}
