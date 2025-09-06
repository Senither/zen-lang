package parser

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/tokens"
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []ParserError

	filePath string

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

type ParserError struct {
	Message  string
	FilePath string
	Token    tokens.Token
}

func (e *ParserError) String() string {
	path := ""
	if e.FilePath != "" {
		path = e.FilePath + ":"
	}

	return fmt.Sprintf(
		"Parser error: %s\n  Token: %q\n  File:  %s%d:%d",
		e.Message, e.Token.Literal, path, e.Token.Line, e.Token.Column,
	)
}

func New(lexer *lexer.Lexer, filePath interface{}) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []ParserError{},
	}

	if path, ok := filePath.(string); ok {
		p.filePath = path
	}

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.INT, p.parseIntegerLiteral)
	p.registerPrefix(tokens.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(tokens.STRING, p.parseStringLiteral)
	p.registerPrefix(tokens.BANG, p.parsePrefixExpression)
	p.registerPrefix(tokens.MINUS, p.parsePrefixExpression)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(tokens.LBRACE, p.parseHashLiteral)
	p.registerPrefix(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(tokens.IF, p.parseIfExpression)
	p.registerPrefix(tokens.WHILE, p.parseWhileExpression)
	p.registerPrefix(tokens.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)
	p.registerInfix(tokens.PLUS, p.parseInfixExpression)
	p.registerInfix(tokens.MINUS, p.parseInfixExpression)
	p.registerInfix(tokens.SLASH, p.parseInfixExpression)
	p.registerInfix(tokens.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tokens.MOD, p.parseInfixExpression)
	p.registerInfix(tokens.EQ, p.parseInfixExpression)
	p.registerInfix(tokens.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.LT, p.parseInfixExpression)
	p.registerInfix(tokens.LT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.GT, p.parseInfixExpression)
	p.registerInfix(tokens.GT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.LPAREN, p.parseCallExpression)
	p.registerInfix(tokens.LBRACKET, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(tokens.EOF) {
		stmt := p.parseStatement()
		if _, ok := stmt.(*ast.EmptyStatement); ok {
			continue
		}

		program.Statements = append(program.Statements, stmt)

		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tokens.VARIABLE:
		return p.parseVariableStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.IMPORT:
		return p.parseImportStatement()
	case tokens.EXPORT:
		return p.parseExportStatement()
	case tokens.BREAK_LOOP:
		return p.parseBreakStatement()
	case tokens.COMMENT:
		return p.parseCommentStatement()
	case tokens.BLOCK_COMMENT_START:
		return p.parseBlockCommentStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t tokens.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected next token to be %q, got %q instead", t, p.peekToken.Type)
	p.errors = append(p.errors, ParserError{
		Message:  msg,
		FilePath: p.filePath,
		Token:    p.peekToken,
	})
}
