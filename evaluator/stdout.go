package evaluator

import (
	"fmt"

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
