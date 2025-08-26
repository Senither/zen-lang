package lexer

import (
	"testing"

	"github.com/senither/zen-lang/tokens"
)

func TestNextToken(t *testing.T) {
	input := `
		var five = 5;
		var pie = 3.14;

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

		"one-word";
		"multiple words";
		'one-word';
		'multiple words';

		=+-!*%/<>
		== !=;
		<= >=;

		@;
	`

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		// Variable assignment to five
		{tokens.VARIABLE, "var"},
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
		// Illegal
		{tokens.ILLEGAL, "@"},
		{tokens.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
