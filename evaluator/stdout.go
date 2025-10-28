package evaluator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/objects"
)

type StandardOut struct {
	Messages []string
	muted    bool
}

var Stdout = &StandardOut{
	Messages: []string{},
	muted:    false,
}

func (s *StandardOut) Write(message string) {
	if !s.muted {
		fmt.Print(message)
	}

	s.Messages = append(s.Messages, message)
}

func (s *StandardOut) ReadAll() []string {
	return s.Messages
}

func (s *StandardOut) Clear() {
	s.Messages = []string{}
}

func (s *StandardOut) Mute(fn func() objects.Object) objects.Object {
	s.muted = true
	rs := fn()
	s.muted = false
	return rs
}

func captureStdoutForBuiltin(
	node *ast.CallExpression,
	fn *objects.ASTAwareBuiltin,
	args []objects.Object,
	env *objects.Environment,
) objects.Object {
	var buf bytes.Buffer

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rs := fn.Fn(node, env, args...)

	w.Close()
	buf.ReadFrom(r)
	os.Stdout = originalStdout

	output := buf.String()
	if output != "" && output != "\n" {
		Stdout.Write(output)
	}

	return rs
}
