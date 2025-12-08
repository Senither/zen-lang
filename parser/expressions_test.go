package parser

import (
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/lexer"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not 'foobar', got %q", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() is not 'foobar', got %q", ident.TokenLiteral())
	}
}

func TestNullLiteralExpression(t *testing.T) {
	input := "null"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NullLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.NullLiteral, got %T", stmt.Expression)
	}

	if literal.TokenLiteral() != "null" {
		t.Errorf("literal.TokenLiteral() is not 'null', got %q", literal.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value is not 5, got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not '5', got %q", literal.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "3.14"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FloatLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 3.14 {
		t.Errorf("literal.Value is not 3.14, got %f", literal.Value)
	}

	if literal.TokenLiteral() != "3.14" {
		t.Errorf("literal.TokenLiteral() is not '3.14', got %q", literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("statement expression is not *ast.StringLiteral. got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got %q", "hello world", literal.Value)
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("statement expression is not *ast.Boolean. got %T", stmt.Expression)
		}

		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got %t", tt.expectedBoolean, boolean.Value)
		}
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	input := "[1, 2 + 3, 4 * 5]"

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayLiteral. got %T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Errorf("array.Elements does not contain 3 elements. got %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "+", 3)
	testInfixExpression(t, array.Elements[2], 4, "*", 5)
}

func TestParsingHashLiteralWithNoKeys(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs is not empty, got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralWithStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs is not %d, got %d", len(expected), len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got %T", key)
			continue
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.HashLiteral, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != len(tests) {
		t.Errorf("hash.Pairs is not %d, got %d", len(tests), len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got %T", key)
			continue
		}

		testFunc, ok := tests[literal.Value]
		if !ok {
			t.Errorf("no test function for key %q found", literal.Value)
			continue
		}

		testFunc(value)
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 2]"

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IndexExpression, got %T", stmt.Expression)
	}

	if indexExp.Left.String() != "myArray" {
		t.Errorf("left side of index expression is not 'myArray', got %q", indexExp.Left.String())
	}

	testInfixExpression(t, indexExp.Index, 1, "+", 2)
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s', got %q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			t.Errorf("exp.Right is not '%d', got %q", tt.value, exp.Right.String())
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 ^ 5", 5, "^", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right) {
			return
		}
	}
}

func TestParsingCompoundInfixAssignments(t *testing.T) {
	tests := []struct {
		input    string
		ident    string
		operator string
		right    int64
	}{
		{"i += 42", "i", "+", 42},
		{"i -= 42", "i", "-", 42},
		{"i *= 42", "i", "*", 42},
		{"i /= 42", "i", "/", 42},
		{"i %= 42", "i", "%", 42},
		{"i ^= 42", "i", "^", 42},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ExpressionStatement, got %T", program.Statements[0])
		}

		assignExp, ok := stmt.Expression.(*ast.AssignmentExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.AssignmentExpression, got %T", stmt.Expression)
		}

		if !testIdentifier(t, assignExp.Left, tt.ident) {
			return
		}

		infixExp, ok := assignExp.Right.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("assignExp.Right is not ast.InfixExpression, got %T", assignExp.Right)
		}

		if !testIdentifier(t, infixExp.Left, tt.ident) {
			return
		}

		if infixExp.Operator != tt.operator {
			t.Errorf("infixExp.Operator is not '%s', got %q", tt.operator, infixExp.Operator)
		}

		if !testIntegerLiteral(t, infixExp.Right, tt.right) {
			return
		}
	}
}

func TestParsingSuffixExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  interface{}
		op    string
	}{
		{"i++", "i", "++"},
		{"i--", "i", "--"},
		{"longerVarName++", "longerVarName", "++"},
		{"longerVarName--", "longerVarName", "--"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.SuffixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.SuffixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.op {
			t.Errorf("exp.Operator is not '%s', got %q", tt.op, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Left, tt.left) {
			t.Errorf("exp.Left is not '%q', got %q", tt.left, exp.Left.String())
		}
	}
}

func TestParsingSuffixExpressionsInBodyDoesntEscapeScope(t *testing.T) {
	input := `func() { i++ } print(i);`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral, got %T", stmt.Expression)
	}

	if len(funcLiteral.Body.Statements) != 1 {
		t.Fatalf("function body does not contain 1 statement, got %d", len(funcLiteral.Body.Statements))
	}

	suffixStmt, ok := funcLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement, got %T", funcLiteral.Body.Statements[0])
	}

	suffixExp, ok := suffixStmt.Expression.(*ast.SuffixExpression)
	if !ok {
		t.Fatalf("function body statement is not ast.SuffixExpression, got %T", suffixStmt.Expression)
	}

	if suffixExp.Operator != "++" {
		t.Errorf("suffixExp.Operator is not '++', got %q", suffixExp.Operator)
	}

	if !testIdentifier(t, suffixExp.Left, "i") {
		return
	}

	printStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement, got %T", program.Statements[1])
	}

	printCall, ok := printStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("printStmt.Expression is not ast.CallExpression, got %T", printStmt.Expression)
	}

	if !testIdentifier(t, printCall.Function, "print") {
		return
	}

	if len(printCall.Arguments) != 1 {
		t.Fatalf("printCall.Arguments does not contain 1 argument, got %d", len(printCall.Arguments))
	}

	if !testIdentifier(t, printCall.Arguments[0], "i") {
		return
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[1], b[0], 2 * [3, 4][1])",
			"add((a * (b[1])), (b[0]), (2 * ([3, 4][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got %q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("alternative is not nil. got %q", exp.Alternative.String())
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got %d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got %d\n", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestIfElseIfElseExpression(t *testing.T) {
	input := []string{
		"if (x < y) { x } elseif (x > y) { y } else { z }",
		"if (x < y) { x } else if (x > y) { y } else { z }",
	}

	for _, expr := range input {
		l := lexer.New(expr)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got %T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Errorf("consequence is not 1 statement. got %d", len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("consequence is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if exp.Intermediary == nil {
			t.Errorf("exp.Intermediaries is nil, expected *ast.IfExpression. got %T\n", exp.Intermediary)
		}

		if !testInfixExpression(t, exp.Intermediary.Condition, "x", ">", "y") {
			return
		}

		if len(exp.Intermediary.Consequence.Statements) != 1 {
			t.Errorf("exp.Intermediaries.Consequence is not 1 statements. got %d\n", len(exp.Intermediary.Consequence.Statements))
		}

		elseifConsequence, ok := exp.Intermediary.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Intermediary.Consequence.Statements[0])
		}

		if !testIdentifier(t, elseifConsequence.Expression, "y") {
			return
		}

		alternative, ok := exp.Intermediary.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Intermediary.Alternative.Statements[0])
		}

		if !testIdentifier(t, alternative.Expression, "z") {
			return
		}
	}
}

func TestWhileExpression(t *testing.T) {
	input := "while (i < 5) { i++ }"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.WhileExpression. got %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "i", "<", 5) {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("while body is not 1 statement. got %d", len(exp.Body.Statements))
	}

	body, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("while body is not ast.ExpressionStatement. got %T", exp.Body.Statements[0])
	}

	if body.String() != "(i++)" {
		t.Errorf("while body is not '(i++)'. got %q", body.String())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "func hello(x, y) { x + y; }"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got %T", stmt.Expression)
	}

	if funcLiteral.Name == nil {
		t.Fatalf("function literal name is nil, expected *ast.Identifier, got %T", funcLiteral.Name)
	}

	testLiteralExpression(t, funcLiteral.Name, "hello")

	if len(funcLiteral.Parameters) != 2 {
		t.Errorf("function literal parameters are not 2. got %d", len(funcLiteral.Parameters))
	}

	testLiteralExpression(t, funcLiteral.Parameters[0], "x")
	testLiteralExpression(t, funcLiteral.Parameters[1], "y")

	if len(funcLiteral.Body.Statements) != 1 {
		t.Errorf("function.Body.Statements does not have 1 statement. got %d", len(funcLiteral.Body.Statements))
	}

	bodyStmt, ok := funcLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got %T", funcLiteral.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestNamedFunctionLiteralParsing(t *testing.T) {
	input := []string{
		`func(x, y) { x + y; }`,
		`func (x, y) { x + y; }`,
	}

	for _, expr := range input {
		l := lexer.New(expr)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got %T", stmt.Expression)
		}

		if len(funcLiteral.Parameters) != 2 {
			t.Errorf("function literal parameters are not 2. got %d", len(funcLiteral.Parameters))
		}

		testLiteralExpression(t, funcLiteral.Parameters[0], "x")
		testLiteralExpression(t, funcLiteral.Parameters[1], "y")

		if len(funcLiteral.Body.Statements) != 1 {
			t.Errorf("function.Body.Statements does not have 1 statement. got %d", len(funcLiteral.Body.Statements))
		}

		bodyStmt, ok := funcLiteral.Body.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got %T", funcLiteral.Body.Statements[0])
		}

		testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
		{input: "func name() {};", expectedParams: []string{}},
		{input: "func name(x) {};", expectedParams: []string{"x"}},
		{input: "func name(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got %T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got %d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	tests := []struct {
		input             string
		expectedFunction  string
		expectedArguments []string
	}{
		{"add(1, 2)", "add", []string{"1", "2"}},
		{"subtract(5, 3)", "subtract", []string{"5", "3"}},
		{"multiply(2, 4)", "multiply", []string{"2", "4"}},
		{"special(1, 2 * 3, 4 + 5)", "special", []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		callExp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got %T", stmt.Expression)
		}

		if !testIdentifier(t, callExp.Function, tt.expectedFunction) {
			return
		}

		if len(callExp.Arguments) != len(tt.expectedArguments) {
			t.Errorf("wrong number of arguments. want %d, got %d", len(tt.expectedArguments), len(callExp.Arguments))
		}

		for i, arg := range tt.expectedArguments {
			if callExp.Arguments[i].String() != arg {
				t.Errorf("argument %d is not %q. got %q", i, arg, callExp.Arguments[i].String())
			}
		}
	}
}

func TestChainExpressionParsing(t *testing.T) {
	tests := []struct {
		input string
		left  string
		right string
	}{
		{"foo.bar", "foo", "bar"},
		{"baz.qux", "baz", "qux"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		chainExp, ok := stmt.Expression.(*ast.ChainExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
		}

		testLiteralExpression(t, chainExp.Left, tt.left)
		testLiteralExpression(t, chainExp.Right, tt.right)
	}
}

func TestNestedChainExpressionParsing(t *testing.T) {
	input := "a.b.c.d"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	chainAExp, ok := stmt.Expression.(*ast.ChainExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
	}

	testLiteralExpression(t, chainAExp.Left, "a")

	chainBExp, ok := chainAExp.Right.(*ast.ChainExpression)
	if !ok {
		t.Fatalf("chainAExp.Right is not ast.ChainExpression. got %T", chainAExp.Right)
	}

	testLiteralExpression(t, chainBExp.Left, "b")

	chainCExp, ok := chainBExp.Right.(*ast.ChainExpression)
	if !ok {
		t.Fatalf("chainBExp.Right is not ast.ChainExpression. got %T", chainBExp.Right)
	}

	testLiteralExpression(t, chainCExp.Left, "c")

	chainDExp, ok := chainCExp.Right.(*ast.Identifier)
	if !ok {
		t.Fatalf("chainCExp.Right is not ast.Identifier. got %T", chainCExp.Right)
	}

	testLiteralExpression(t, chainDExp, "d")
}

func TestChainCallExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		left     string
		function string
		args     []string
	}{
		{"obj.foo(5)", "obj", "foo", []string{"5"}},
		{"obj.bar(10, 20)", "obj", "bar", []string{"10", "20"}},
		{"obj.baz(1 + 2, 3 * 4)", "obj", "baz", []string{"(1 + 2)", "(3 * 4)"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		chainExp, ok := stmt.Expression.(*ast.ChainExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
		}

		callExp, ok := chainExp.Right.(*ast.CallExpression)
		if !ok {
			t.Fatalf("chainExp.Right is not ast.CallExpression. got %T", chainExp.Right)
		}

		testLiteralExpression(t, callExp.Function, tt.function)

		if len(callExp.Arguments) != len(tt.args) {
			t.Errorf("wrong number of arguments. want %d, got %d", len(tt.args), len(callExp.Arguments))
		}

		for i, arg := range tt.args {
			if callExp.Arguments[i].String() != arg {
				t.Errorf("argument %d is not %q. got %q", i, arg, callExp.Arguments[i].String())
			}
		}
	}
}

func TestChainIndexExpressionParsing(t *testing.T) {
	tests := []struct {
		input string
		ident string
		left  string
		index any
	}{
		{"data.items[0]", "data", "items", 0},
		{"config.values[key]", "config", "values", "key"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, nil)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
		}

		chainExp, ok := stmt.Expression.(*ast.ChainExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
		}

		indexExp, ok := chainExp.Right.(*ast.IndexExpression)
		if !ok {
			t.Fatalf("chainExp.Right is not ast.IndexExpression. got %T", chainExp.Right)
		}

		testLiteralExpression(t, chainExp.Left, tt.ident)
		testLiteralExpression(t, indexExp.Left, tt.left)
		testLiteralExpression(t, indexExp.Index, tt.index)
	}
}

func TestChainAssignmentExpressionParsing(t *testing.T) {
	input := "obj.prop = 42"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	chainExp, ok := stmt.Expression.(*ast.ChainExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
	}

	assignExp, ok := chainExp.Right.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("chainExp.Right is not ast.AssignmentExpression. got %T", chainExp.Right)
	}

	objAssignExp, ok := assignExp.Right.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("assignExp.Right is not ast.AssignmentExpression. got %T", assignExp.Right)
	}

	testLiteralExpression(t, assignExp.Left, "obj")
	testLiteralExpression(t, objAssignExp.Left, "prop")
	testLiteralExpression(t, objAssignExp.Right, 42)
}

func TestChainIndexAssignmentExpressionParsing(t *testing.T) {
	input := "data.items[0] = 100"

	l := lexer.New(input)
	p := New(l, nil)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T", program.Statements[0])
	}

	chainExp, ok := stmt.Expression.(*ast.ChainExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ChainExpression. got %T", stmt.Expression)
	}

	assignExp, ok := chainExp.Right.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("chainExp.Right is not ast.AssignmentExpression. got %T", chainExp.Right)
	}

	indexExp, ok := assignExp.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("assignExp.Left is not ast.IndexExpression. got %T", assignExp.Left)
	}

	testLiteralExpression(t, chainExp.Left, "data")
	testLiteralExpression(t, indexExp.Left, "items")
	testLiteralExpression(t, indexExp.Index, 0)
	testLiteralExpression(t, assignExp.Right, 100)
}
