package lexer

import (
	"strings"

	"github.com/senither/zen-lang/tokens"
)

type Lexer struct {
	input         string
	position      int
	readPosition  int
	currentLine   int
	currentColumn int
	ch            byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input:         input,
		currentLine:   1,
		currentColumn: 0,
	}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() tokens.Token {
	var token tokens.Token

	l.skipWhitespace()

	switch l.ch {
	case ';':
		token = newToken(tokens.SEMICOLON, l)
	case '"':
		token = newTokenWithValue(tokens.STRING, l, l.readString('"'))
	case '\'':
		token = newTokenWithValue(tokens.STRING, l, l.readString('\''))

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.EQ, l, string(ch)+string(l.ch))
		} else {
			token = newToken(tokens.ASSIGN, l)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.NOT_EQ, l, string(ch)+string(l.ch))
		} else {
			token = newToken(tokens.BANG, l)
		}
	case '/':
		switch l.peekChar() {
		case '/':
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.COMMENT, l, string(ch)+string(l.ch))
		case '*':
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.BLOCK_COMMENT_START, l, string(ch)+string(l.ch))
		default:
			token = newToken(tokens.SLASH, l)
		}
	case '+':
		token = newToken(tokens.PLUS, l)
	case '-':
		token = newToken(tokens.MINUS, l)
	case '*':
		switch l.peekChar() {
		case '/':
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.BLOCK_COMMENT_END, l, string(ch)+string(l.ch))
		default:
			token = newToken(tokens.ASTERISK, l)
		}
	case '%':
		token = newToken(tokens.MOD, l)
	case '>':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.GT_EQ, l, string(ch)+string(l.ch))
		default:
			token = newToken(tokens.GT, l)
		}
	case '<':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			token = newTokenWithValue(tokens.LT_EQ, l, string(ch)+string(l.ch))
		default:
			token = newToken(tokens.LT, l)
		}

	case ',':
		token = newToken(tokens.COMMA, l)
	case ':':
		token = newToken(tokens.COLON, l)

	case '(':
		token = newToken(tokens.LPAREN, l)
	case ')':
		token = newToken(tokens.RPAREN, l)
	case '{':
		token = newToken(tokens.LBRACE, l)
	case '}':
		token = newToken(tokens.RBRACE, l)
	case '[':
		token = newToken(tokens.LBRACKET, l)
	case ']':
		token = newToken(tokens.RBRACKET, l)

	case 0:
		token = newTokenWithValue(tokens.EOF, l, "")

	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()

			if literal == "else" {
				token := l.readIfElseToken()
				if token != nil {
					return *token
				}
			}

			return newTokenWithValue(tokens.LookupIdent(literal), l, literal)
		} else if isDigit(l.ch) {
			return l.readNumberToken()
		} else {
			token = newToken(tokens.ILLEGAL, l)
		}
	}

	l.readChar()

	return token
}

func (l *Lexer) readIfElseToken() *tokens.Token {
	l.skipWhitespace()

	if l.ch == 0 {
		return nil
	}

	if !isLetter(l.ch) {
		return nil
	}

	start := l.position
	nextIdent := l.readIdentifier()

	if nextIdent != "if" {
		// Reset the position to fix the lexer state
		l.position = start
		l.readPosition = start + 1
		l.ch = l.input[l.position]

		return nil
	}

	token := newTokenWithValue(tokens.ELSE_IF, l, "else if")
	return &token
}

func newTokenWithValue(tokenType tokens.TokenType, l *Lexer, value string) tokens.Token {
	column := l.currentColumn
	if len(value) > 1 {
		column = l.currentColumn - len(value)
	} else if isDynamicToken(tokenType) {
		column--
	}

	return tokens.Token{
		Type:    tokenType,
		Literal: value,
		Column:  column,
		Line:    l.currentLine,
	}
}

func newToken(tokenType tokens.TokenType, l *Lexer) tokens.Token {
	return tokens.Token{
		Type:    tokenType,
		Literal: string(l.ch),
		Column:  l.currentColumn,
		Line:    l.currentLine,
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString(endChar byte) string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == endChar || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
	l.currentColumn++
}

func (l *Lexer) readNumberToken() tokens.Token {
	val := l.readNumber()
	if strings.Contains(val, ".") {
		return newTokenWithValue(tokens.FLOAT, l, val)
	}

	return newTokenWithValue(tokens.INT, l, val)
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.currentLine++
			l.currentColumn = 0
		}

		l.readChar()
	}
}

func isDynamicToken(token tokens.TokenType) bool {
	return token == tokens.IDENT || token == tokens.INT
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
