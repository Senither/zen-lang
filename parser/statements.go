package parser

import (
	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.curToken, Mutable: false}

	if p.peekTokenIs(tokens.MUTABLE) {
		stmt.Mutable = true
		p.nextToken()
	}

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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(tokens.STRING) {
		return nil
	}

	stmt.Path = p.curToken.Literal

	if p.peekTokenIs(tokens.IMPORT_ALIAS) {
		p.nextToken()
		p.nextToken()

		stmt.Aliased = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExportStatement() *ast.ExportStatement {
	stmt := &ast.ExportStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if !p.expectPeek(tokens.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

	if !p.expectPeek(tokens.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseCommentStatement() *ast.EmptyStatement {
	curLine := p.curToken.Line

	for curLine == p.curToken.Line && !p.curTokenIs(tokens.EOF) {
		p.nextToken()
	}

	return &ast.EmptyStatement{Token: p.curToken}
}

func (p *Parser) parseBlockCommentStatement() *ast.EmptyStatement {
	for !p.curTokenIs(tokens.BLOCK_COMMENT_END) && !p.curTokenIs(tokens.EOF) {
		p.nextToken()
	}

	p.nextToken()

	return &ast.EmptyStatement{Token: p.curToken}
}
