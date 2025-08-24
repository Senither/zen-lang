package tokens

type Token struct {
	Type    TokenType
	Literal string
	Column  int
	Line    int
}

type TokenType string

const (
	ILLEGAL   TokenType = "ILLEGAL" // Illegal token
	EOF       TokenType = "EOF"     // End of file
	SEMICOLON TokenType = ";"

	// Identifiers + literals
	VARIABLE TokenType = "VARIABLE" // Variable identifiers
	IDENT    TokenType = "IDENT"    // add, foobar, x, y, ...
	INT      TokenType = "INT"      // 1343456
	STRING   TokenType = "STRING"   // "string"

	// String literals
	DOUBLE_QUOTE TokenType = "\"" // "string"
	SINGLE_QUOTE TokenType = "'"  // 'string'

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MOD      TokenType = "%"
	GT       TokenType = ">"
	LT       TokenType = "<"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// Delimiters
	COMMA TokenType = ","

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Comments
	COMMENT             TokenType = "COMMENT"
	BLOCK_COMMENT_START TokenType = "BLOCK_COMMENT_START"
	BLOCK_COMMENT_END   TokenType = "BLOCK_COMMENT_END"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	ELSE_IF  TokenType = "ELSE_IF"
	RETURN   TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"var":     VARIABLE,
	"func":    FUNCTION,
	"fn":      FUNCTION,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"elseif":  ELSE_IF,
	"else if": ELSE_IF,
	"return":  RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
