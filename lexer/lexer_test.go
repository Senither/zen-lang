package lexer

import (
	"testing"

	"github.com/senither/zen-lang/tokens"
)

func TestNextToken(t *testing.T) {
	input := `
		var mut five = 5;
		var pie = 3.14;
		var val = 9f;
		var nil = null;

		var add = func(x, y) {
			return x + y;
		}

		func multiply(x, y) {
			return x * y;
		}

		var result = add(five, 10);

		if (5 < 10) {
			return true;
		} elseif (5 > 10) {
			return false;
		} else {
			return "Something went wrong";
		}

		while (i < 10) {
			i++;
			break;
		}

		"one-word";
		"multiple words";
		'one-word';
		'multiple words';

		=+-!*^%/<>
		== !=;
		<= >=;
		&& ||;

		++ --;

		[1, 2];
		{"foo": "bar"};
		obj.foo(5);

		@;
	`

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		// Variable assignment to five
		{tokens.VARIABLE, "var"},
		{tokens.MUTABLE, "mut"},
		{tokens.IDENT, "five"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "5"},
		{tokens.SEMICOLON, ";"},
		// Variable assignment to pie
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "pie"},
		{tokens.ASSIGN, "="},
		{tokens.FLOAT, "3.14"},
		{tokens.SEMICOLON, ";"},
		// Variable assignment to val
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "val"},
		{tokens.ASSIGN, "="},
		{tokens.FLOAT, "9"},
		{tokens.SEMICOLON, ";"},
		// Variable assignment to nil
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "nil"},
		{tokens.ASSIGN, "="},
		{tokens.NULL, "null"},
		{tokens.SEMICOLON, ";"},
		// Variable for add function
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "add"},
		{tokens.ASSIGN, "="},
		{tokens.FUNCTION, "func"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "y"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.IDENT, "x"},
		{tokens.PLUS, "+"},
		{tokens.IDENT, "y"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		// Named function for multiply function
		{tokens.FUNCTION, "func"},
		{tokens.IDENT, "multiply"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "y"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.IDENT, "x"},
		{tokens.ASTERISK, "*"},
		{tokens.IDENT, "y"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		// Variable assignment to result for function call
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "result"},
		{tokens.ASSIGN, "="},
		{tokens.IDENT, "add"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "five"},
		{tokens.COMMA, ","},
		{tokens.INT, "10"},
		{tokens.RPAREN, ")"},
		{tokens.SEMICOLON, ";"},
		// If & if-else statement
		{tokens.IF, "if"},
		{tokens.LPAREN, "("},
		{tokens.INT, "5"},
		{tokens.LT, "<"},
		{tokens.INT, "10"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.TRUE, "true"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.ELSE_IF, "elseif"},
		{tokens.LPAREN, "("},
		{tokens.INT, "5"},
		{tokens.GT, ">"},
		{tokens.INT, "10"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.FALSE, "false"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.ELSE, "else"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.STRING, "Something went wrong"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		// While loops
		{tokens.WHILE, "while"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "i"},
		{tokens.LT, "<"},
		{tokens.INT, "10"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.IDENT, "i"},
		{tokens.INCREMENT, "++"},
		{tokens.SEMICOLON, ";"},
		{tokens.BREAK_LOOP, "break"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		// String literals
		{tokens.STRING, "one-word"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, "multiple words"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, "one-word"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, "multiple words"},
		{tokens.SEMICOLON, ";"},
		// Expression operators
		{tokens.ASSIGN, "="},
		{tokens.PLUS, "+"},
		{tokens.MINUS, "-"},
		{tokens.BANG, "!"},
		{tokens.ASTERISK, "*"},
		{tokens.CARET, "^"},
		{tokens.MOD, "%"},
		{tokens.SLASH, "/"},
		{tokens.LT, "<"},
		{tokens.GT, ">"},
		{tokens.EQ, "=="},
		{tokens.NOT_EQ, "!="},
		{tokens.SEMICOLON, ";"},
		{tokens.LT_EQ, "<="},
		{tokens.GT_EQ, ">="},
		{tokens.SEMICOLON, ";"},
		{tokens.AND, "&&"},
		{tokens.OR, "||"},
		{tokens.SEMICOLON, ";"},
		// Increment & Decrement
		{tokens.INCREMENT, "++"},
		{tokens.DECREMENT, "--"},
		{tokens.SEMICOLON, ";"},
		// Array literals
		{tokens.LBRACKET, "["},
		{tokens.INT, "1"},
		{tokens.COMMA, ","},
		{tokens.INT, "2"},
		{tokens.RBRACKET, "]"},
		{tokens.SEMICOLON, ";"},
		// HashMap literals
		{tokens.LBRACE, "{"},
		{tokens.STRING, "foo"},
		{tokens.COLON, ":"},
		{tokens.STRING, "bar"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		// Chained call expression on HashMap
		{tokens.IDENT, "obj"},
		{tokens.PERIOD, "."},
		{tokens.IDENT, "foo"},
		{tokens.LPAREN, "("},
		{tokens.INT, "5"},
		{tokens.RPAREN, ")"},
		{tokens.SEMICOLON, ";"},
		// Illegal
		{tokens.ILLEGAL, "@"},
		{tokens.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong.\nexpected %q\ngot %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong.\nexpected %q\ngot %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	tok := l.NextToken()
	if tok.Type != tokens.EOF {
		t.Fatalf("expected EOF token at end of input, got %q (value: %q)", tok.Type, tok.Literal)
	}
}

func TestNextTokenEscapedString(t *testing.T) {
	input := `
		"\\ \ \'" \\ \ \''
		"Hello\"World" \""
		"Hello\nWorld" \n
		"Hello\rWorld" \r
		"Hello\tWorld" \t
		"Hello\bWorld" \b
		"Hello\aWorld" \a
		"Hello\0World" \0
	`

	l := New(input)

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		// Backslash escape
		{tokens.STRING, `\ \ '`},
		{tokens.ILLEGAL, `\`},
		{tokens.ILLEGAL, `\`},
		{tokens.ILLEGAL, `\`},
		{tokens.ILLEGAL, `\`},
		{tokens.STRING, ""},
		// Quote escape
		{tokens.STRING, `Hello"World`},
		{tokens.ILLEGAL, `\`},
		{tokens.STRING, ``},
		// New line escape
		{tokens.STRING, "Hello\nWorld"},
		{tokens.ILLEGAL, "\\"},
		{tokens.IDENT, "n"},
		// Reset escape
		{tokens.STRING, "Hello\rWorld"},
		{tokens.ILLEGAL, "\\"},
		{tokens.IDENT, "r"},
		// Tab escape
		{tokens.STRING, "Hello\tWorld"},
		{tokens.ILLEGAL, "\\"},
		{tokens.IDENT, "t"},
		// Backspace escape
		{tokens.STRING, "Hello\bWorld"},
		{tokens.ILLEGAL, "\\"},
		{tokens.IDENT, "b"},
		// Alert escape
		{tokens.STRING, "Hello\aWorld"},
		{tokens.ILLEGAL, "\\"},
		{tokens.IDENT, "a"},
		// Null byte escape
		{tokens.STRING, "Hello\x00World"},
		{tokens.ILLEGAL, "\\"},
		{tokens.INT, "0"},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong.\nexpected %q,\ngot %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong.\nexpected %q,\ngot %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	tok := l.NextToken()
	if tok.Type != tokens.EOF {
		t.Fatalf("expected EOF token at end of input, got %q (value: %q)", tok.Type, tok.Literal)
	}
}
