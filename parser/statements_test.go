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
	p := New(l, nil)

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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got %T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestImportStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedFile  string
		expectedAlias interface{}
	}{
		{"import 'file'", "file", nil},
		{"import 'file' as f", "file", "f"},
		{"import \"another_file\" as af", "another_file", "af"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		importStmt, ok := stmt.(*ast.ImportStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ImportStatement. got %T", stmt)
		}

		if importStmt.TokenLiteral() != "import" {
			t.Fatalf("importStmt.TokenLiteral is not 'import', got %q", importStmt.TokenLiteral())
		}

		if importStmt.Path != tt.expectedFile {
			t.Errorf("importStmt.ImportPath is not %q. got %q", tt.expectedFile, importStmt.Path)
		}

		if importStmt.Aliased != nil && importStmt.Aliased.Value != tt.expectedAlias {
			t.Errorf("importStmt.Aliased is not %q. got %q", tt.expectedAlias, importStmt.Aliased.Value)
		}
	}
}

func TestExportStatements(t *testing.T) {
	tests := []struct {
		input       string
		expectedAst string
	}{
		{"export someValue", "someValue"},
		{"export someValue;", "someValue"},
		{"export func(x) {x}", "func (x) { x }"},
		{"export func(x) {x};", "func (x) { x }"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		exportStmt, ok := stmt.(*ast.ExportStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExportStatement. got %T", stmt)
		}

		if exportStmt.TokenLiteral() != "export" {
			t.Fatalf("exportStmt.TokenLiteral is not 'export', got %q", exportStmt.TokenLiteral())
		}

		if exportStmt.Value.String() != tt.expectedAst {
			t.Errorf("exportStmt.ExportValue is not %q. got %q", tt.expectedAst, exportStmt.Value.String())
		}
	}
}

func TestCommentStatements(t *testing.T) {
	input := `
		// This is a comment
		var x = 5;
		// Another comment
		var y = 10;
	`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got %d", len(program.Statements))
	}

	firstStmt, ok := program.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VariableStatement. got %T", program.Statements[0])
	}

	if firstStmt.Name.Value != "x" {
		t.Errorf("varStmt.Name.Value is not 'x'. got %q", firstStmt.Name.Value)
	}

	secondStmt, ok := program.Statements[1].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.VariableStatement. got %T", program.Statements[1])
	}

	if secondStmt.Name.Value != "y" {
		t.Errorf("varStmt.Name.Value is not 'y'. got %q", secondStmt.Name.Value)
	}
}

func TestBlockCommentStatements(t *testing.T) {
	input := `
		/* This is a block comment */
		var x = 5;
		/*
		This is a longer block comment
		That spans multiple lines
		*/
		var y = 10;
	`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got %d", len(program.Statements))
	}

	firstStmt, ok := program.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VariableStatement. got %T", program.Statements[0])
	}

	if firstStmt.Name.Value != "x" {
		t.Errorf("varStmt.Name.Value is not 'x'. got %q", firstStmt.Name.Value)
	}

	secondStmt, ok := program.Statements[1].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.VariableStatement. got %T", program.Statements[1])
	}

	if secondStmt.Name.Value != "y" {
		t.Errorf("varStmt.Name.Value is not 'y'. got %q", secondStmt.Name.Value)
	}
}

func TestBreakLoopStatement(t *testing.T) {
	input := `break;`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	breakStmt, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BreakStatement. got %T", program.Statements[0])
	}

	if breakStmt.TokenLiteral() != "break" {
		t.Fatalf("breakStmt.TokenLiteral is not 'break', got %q", breakStmt.TokenLiteral())
	}
}

func TestContinueLoopStatement(t *testing.T) {
	input := `continue;`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	continueStmt, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ContinueStatement. got %T", program.Statements[0])
	}

	if continueStmt.TokenLiteral() != "continue" {
		t.Fatalf("continueStmt.TokenLiteral is not 'continue', got %q", continueStmt.TokenLiteral())
	}
}
