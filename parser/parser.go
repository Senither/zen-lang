package parser

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/tokens"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  tokens.Token
	peekToken tokens.Token
	errors    []ParserError
}

type ParserError struct {
	Message string
	Token   tokens.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []ParserError{},
	}

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
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tokens.VARIABLE:
		return p.parseVariableStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
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
	p.errors = append(p.errors, ParserError{Message: msg, Token: p.peekToken})
}
