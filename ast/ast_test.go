package ast

import (
	"testing"

	"github.com/senither/zen-lang/tokens"
)

func TestString(t *testing.T) {
	expected := "var myVar = anotherVar;"
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: tokens.Token{Type: tokens.VARIABLE, Literal: "var"},
				Name: &Identifier{
					Token: tokens.Token{
						Type:    tokens.IDENT,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tokens.Token{Type: tokens.VARIABLE, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != expected {
		t.Errorf("program.String() is not %q. got=%q", expected, program.String())
	}
}
