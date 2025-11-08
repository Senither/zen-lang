package parser

import (
	"fmt"
	"strconv"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // < or >
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      //-x or !y
	CALL        // myFunction(x)
	INDEX       // array[index]
)

var precedences = map[tokens.TokenType]int{
	tokens.EQ:       EQUALS,
	tokens.NOT_EQ:   EQUALS,
	tokens.LT:       LESSGREATER,
	tokens.GT:       LESSGREATER,
	tokens.LT_EQ:    LESSGREATER,
	tokens.GT_EQ:    LESSGREATER,
	tokens.PLUS:     SUM,
	tokens.MINUS:    SUM,
	tokens.SLASH:    PRODUCT,
	tokens.ASTERISK: PRODUCT,
	tokens.CARET:    PRODUCT,
	tokens.MOD:      PRODUCT,
	tokens.LPAREN:   CALL,
	tokens.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(tokens.ASSIGN) {
		p.nextToken()

		return p.parseAssignmentExpression(ident)
	} else if p.peekTokenIs(tokens.INCREMENT) || p.peekTokenIs(tokens.DECREMENT) {
		p.nextToken()

		return p.parseSuffixExpression(ident)
	} else if p.peekTokenIs(tokens.PERIOD) {
		p.nextToken()

		return p.parseChainExpression(ident)
	}

	return ident
}

func (p *Parser) parseChainExpression(left ast.Expression) ast.Expression {
	chain := &ast.ChainExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp := p.parseExpression(LOWEST)

	switch exp.(type) {
	case *ast.Identifier, *ast.CallExpression, *ast.ChainExpression, *ast.IndexExpression:
		chain.Right = exp

	case *ast.AssignmentExpression:
		chain.Right = &ast.AssignmentExpression{
			Token: p.curToken,
			Left:  chain.Left,
			Right: exp,
		}

	default:
		msg := fmt.Sprintf("unexpected chained expression, got %T", exp)
		p.errors = append(p.errors, ParserError{
			Message:  msg,
			FilePath: p.filePath,
			Token:    p.curToken,
		})

		return nil
	}

	return chain
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Right = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{
			Message:  msg,
			FilePath: p.filePath,
			Token:    p.curToken,
		})
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	literal := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, ParserError{
			Message:  msg,
			FilePath: p.filePath,
			Token:    p.curToken,
		})
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(tokens.TRUE)}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(tokens.RBRACKET)

	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(tokens.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)
		if !p.expectPeek(tokens.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(tokens.RBRACE) && !p.expectPeek(tokens.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(tokens.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParserError{
		Message:  msg,
		FilePath: p.filePath,
		Token:    p.curToken,
	})
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseSuffixExpression(left ast.Expression) ast.Expression {
	expression := &ast.SuffixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token:        p.curToken,
		Intermediary: nil,
	}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	for p.peekTokenIs(tokens.ELSE_IF) {
		p.nextToken()
		expression.Intermediary = p.parseIfExpression().(*ast.IfExpression)
	}

	if p.peekTokenIs(tokens.ELSE) {
		p.nextToken()

		if !p.expectPeek(tokens.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	funcLiteral := &ast.FunctionLiteral{
		Token: p.curToken,
		Name:  nil,
	}

	if p.peekTokenIs(tokens.IDENT) {
		p.nextToken()
		funcLiteral.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	funcLiteral.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	funcLiteral.Body = p.parseBlockStatement()

	return funcLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(tokens.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(tokens.RBRACE) && !p.curTokenIs(tokens.EOF) {
		block.Statements = append(block.Statements, p.parseStatement())
		p.nextToken()
	}

	return block
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(tokens.RPAREN)

	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RBRACKET) {
		return nil
	}

	if !p.peekTokenIs(tokens.ASSIGN) {
		return exp
	}

	p.nextToken()

	return p.parseAssignmentExpression(exp)
}

func (p *Parser) parseExpressionList(end tokens.TokenType) []ast.Expression {
	expressions := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return expressions
	}

	p.nextToken()
	expressions = append(expressions, p.parseExpression(LOWEST))

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return expressions
}
