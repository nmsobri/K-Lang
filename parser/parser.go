package parser

import (
	"Klang/ast"
	"Klang/lexer"
	"Klang/token"
	"strconv"
)

const (
	_ = iota
	LOWEST
	SUM
	PRODUCT
)

var precedence = map[token.TokenType]int{
	token.PLUS:  SUM,
	token.MINUS: SUM,
	token.STAR:  PRODUCT,
	token.SLASH: PRODUCT,
}

type (
	prefixFunc func() ast.Expression
	infixFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	Lexer        *lexer.Lexer
	CurrentToken token.Token
	PeekToken    token.Token
	prefixFunc   map[token.TokenType]prefixFunc
	infixFunc    map[token.TokenType]infixFunc
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		Lexer: lex,
	}

	p.prefixFunc = make(map[token.TokenType]prefixFunc)
	p.infixFunc = make(map[token.TokenType]infixFunc)

	// prefix function
	p.registerPrefixFunction(token.INTEGER, p.parseInteger)
	p.registerPrefixFunction(token.FLOATING, p.parseFloating)
	p.registerPrefixFunction(token.TRUE, p.parseBoolean)
	p.registerPrefixFunction(token.FALSE, p.parseBoolean)
	p.registerPrefixFunction(token.IF, p.parseIfExpression)

	// infix function
	p.registerInfixFunction(token.PLUS, p.parseInfixExpression)
	p.registerInfixFunction(token.MINUS, p.parseInfixExpression)
	p.registerInfixFunction(token.STAR, p.parseInfixExpression)
	p.registerInfixFunction(token.SLASH, p.parseInfixExpression)

	// prime the tokens
	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) NextToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.Lexer.NextToken()
}

func (p *Parser) expectPeek(tokType token.TokenType) bool {
	if p.PeekToken.Type == tokType {
		p.NextToken()
		return true
	}

	return false
}

func (p *Parser) peekTokenIs(tokType token.TokenType) bool {
	return p.PeekToken.Type == tokType
}

func (p *Parser) currentTokenIs(tokType token.TokenType) bool {
	return p.CurrentToken.Type == tokType
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.CurrentToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.parseLetStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	letStmt := &ast.LetStatement{Token: p.CurrentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	letStmt.Name = p.parseIdentifier().(*ast.Identifier)

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.NextToken() // advance to value expression
	letStmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken() // consume semicolon if peek token is semicolon, since its optional on repl
	}

	return letStmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	returnStmt := &ast.ReturnStatement{Token: p.CurrentToken}

	p.NextToken() // advance to value expression
	returnStmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return returnStmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expressionStmt := &ast.ExpressionStatement{Token: p.CurrentToken}
	expressionStmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return expressionStmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	return ident
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.getPrefixFunction(p.CurrentToken)

	if prefix == nil {
		return nil
	}

	// 1 + 2; #pass
	// 1 +;

	left := prefix()

	// pratt parser
	// keep parsing, if next token is more important than the current one
	// else process what we have so far
	for p.peekTokenPrecedence() > precedence {
		infix := p.getInfixFunction(p.PeekToken)

		if infix == nil {
			return left
		}

		p.NextToken() // consume the infix operator
		left = infix(left)
	}

	return left
}

func (p *Parser) parseInteger() ast.Expression {
	val, err := strconv.ParseInt(p.CurrentToken.Literal, 10, 64)

	if err != nil {
		val = 0
	}

	integer := &ast.IntegerLiteral{Token: p.CurrentToken, Value: val}
	return integer
}

func (p *Parser) parseFloating() ast.Expression {
	val, err := strconv.ParseFloat(p.CurrentToken.Literal, 64)

	if err != nil {
		val = 0
	}

	float := &ast.FloatLiteral{Token: p.CurrentToken, Value: val}
	return float
}

func (p *Parser) getPrefixFunction(tok token.Token) prefixFunc {
	if fn, ok := p.prefixFunc[tok.Type]; ok {
		return fn
	}

	return nil
}

func (p *Parser) getInfixFunction(tok token.Token) infixFunc {
	if fn, ok := p.infixFunc[tok.Type]; ok {
		return fn
	}

	return nil
}

func (p *Parser) registerPrefixFunction(tokType token.TokenType, fn prefixFunc) {
	p.prefixFunc[tokType] = fn
}

func (p *Parser) registerInfixFunction(tokType token.TokenType, fn infixFunc) {
	p.infixFunc[tokType] = fn
}

func (p *Parser) peekTokenPrecedence() int {
	if prec, ok := precedence[p.PeekToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infix := &ast.InfixExpression{Token: p.CurrentToken, Left: left, Operator: p.CurrentToken.Literal}

	prec := precedence[p.CurrentToken.Type]
	p.NextToken() // advance to the right operand of the operator

	infix.Right = p.parseExpression(prec)
	return infix
}

func (p *Parser) parseBoolean() ast.Expression {
	val := false

	if p.CurrentToken.Literal == "true" {
		val = true
	}

	boolean := &ast.BooleanLiteral{Token: p.CurrentToken, Value: val}
	return boolean
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpr := &ast.IfExpression{Token: p.CurrentToken}

	p.NextToken() // advance to condition expression
	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifExpr.IfArm = p.parseBlockStatement().(*ast.BlockStatement)

	if p.peekTokenIs(token.ELSE) {
		p.NextToken() // consume the `else` token

		// consume the `{`
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		ifExpr.ElseArm = p.parseBlockStatement().(*ast.BlockStatement)
	}

	return ifExpr
}

// TODO continue here
func (p *Parser) parseBlockStatement() ast.Expression {
	block := &ast.BlockStatement{Token: p.CurrentToken}
	block.Statements = []ast.Statement{}

	p.NextToken() // advance to the block body

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.NextToken()
	}

	return block
}
