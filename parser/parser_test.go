package parser

import (
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
		var x = 5;
		var y = 10;
		var hello = "world";
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"hello"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			t.Errorf("TestVarStatements failed for statement %d", i)
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("Expected token literal 'var', got '%s'", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Expected *ast.VarStatement, got %T", s)
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value: Expected identifier %s, got %s", name, varStmt.Name.Value)
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral(): Expected token literal %s, got %v", name, varStmt.Name)
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 'Alexis';
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected *ast.ReturnStatement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

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
