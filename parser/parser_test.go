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

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of exp not handled. got %T", exp)

	return false
}

func testIntegerLiteral(t *testing.T, literal ast.Expression, value int64) bool {
	integer, ok := literal.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integer is not ast.IntegerLiteral. got %T", literal)
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

func testIdentifier(t *testing.T, ident ast.Expression, value string) bool {
	identifier, ok := ident.(*ast.Identifier)
	if !ok {
		t.Errorf("ident is not ast.Identifier. got %T", ident)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %q. got %q", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral is not %q. got %q", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	expression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression. got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, expression.Left, left) {
		return false
	}

	if expression.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, expression.Operator)
		return false
	}

	if !testLiteralExpression(t, expression.Right, right) {
		return false
	}

	return true
}
