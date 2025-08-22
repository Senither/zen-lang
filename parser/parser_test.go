package parser

import (
	"fmt"
	"testing"

	"github.com/senither/zen-lang/ast"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("Parser error: %v @ %d:%d", err.Message, err.Token.Line, err.Token.Column)
	}

	t.FailNow()
}

func testIntegerLiteral(t *testing.T, literal ast.Expression, value int64) bool {
	integer, ok := literal.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not ast.IntegerLiteral. got %T", literal)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d. got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral is not %d. got %q", value, integer.TokenLiteral())
		return false
	}

	return true
}
