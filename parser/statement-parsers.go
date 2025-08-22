package parser

import (
	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.curToken}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we reach the end of the line (;)
	for !p.curTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: We're skipping the expressions until we reach the end of the line (;)
	for !p.curTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
