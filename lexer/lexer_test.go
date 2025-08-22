package lexer

import (
	"testing"

	"github.com/senither/zen-lang/tokens"
)

func TestNextToken(t *testing.T) {
	input := `
		var five = 5;

		var add = fn(x, y) {
			return x + y;
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
		{tokens.EOL, ";"},
		// Variable for add function
		{tokens.VARIABLE, "var"},
		{tokens.IDENT, "add"},
		{tokens.ASSIGN, "="},
		{tokens.FUNCTION, "fn"},
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
		{tokens.EOL, ";"},
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
		{tokens.EOL, ";"},
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
		{tokens.EOL, ";"},
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
		{tokens.EOL, ";"},
		{tokens.RBRACE, "}"},
		{tokens.ELSE, "else"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.STRING, "Something went wrong"},
		{tokens.EOL, ";"},
		{tokens.RBRACE, "}"},
		// String literals
		{tokens.STRING, "one-word"},
		{tokens.EOL, ";"},
		{tokens.STRING, "multiple words"},
		{tokens.EOL, ";"},
		{tokens.STRING, "one-word"},
		{tokens.EOL, ";"},
		{tokens.STRING, "multiple words"},
		{tokens.EOL, ";"},
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
		{tokens.EOL, ";"},
		// Illegal
		{tokens.ILLEGAL, "@"},
		{tokens.EOL, ";"},
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
