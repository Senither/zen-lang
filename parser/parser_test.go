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
