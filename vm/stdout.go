package vm

import (
	"bytes"
	"fmt"
	"os"

	"github.com/senither/zen-lang/objects"
)

type StandardOut struct {
	messages []string
	muted    bool
}

var Stdout = &StandardOut{
	messages: []string{},
	muted:    false,
}

func (s *StandardOut) Write(message string) {
	if !s.muted {
		fmt.Print(message)
	}

	s.messages = append(s.messages, message)
}

func (s *StandardOut) ReadAll() []string {
	return s.messages
}

func (s *StandardOut) Clear() {
	s.messages = []string{}
}

func (s *StandardOut) Mute(fn func() objects.Object) objects.Object {
	s.muted = true
	rs := fn()
	s.muted = false
	return rs
}

func captureStdoutForBuiltin(
	fn *objects.Builtin,
	args []objects.Object,
) objects.Object {
	var buf bytes.Buffer

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rs, err := fn.Fn(args...)

	w.Close()
	buf.ReadFrom(r)
	os.Stdout = originalStdout

	output := buf.String()
	if output != "" && output != "\n" {
		Stdout.Write(output)
	}

	if err != nil {
		return objects.NativeErrorToErrorObject(err)
	}

	return rs
}
