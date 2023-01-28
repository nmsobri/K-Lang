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
	COMPARE
	PREFIX
	CALL
	INDEX
)

// this precedence is only needed by infix operator
// since it is being used by parseExpression to determine
// wether to continue parsing or not
var precedence = map[token.TokenType]int{
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.STAR:          PRODUCT,
	token.SLASH:         PRODUCT,
	token.GREATER:       COMPARE,
	token.GREATER_EQUAL: COMPARE,
	token.LESSER:        COMPARE,
	token.LESSER_EQUAL:  COMPARE,
	token.EQUAL:         COMPARE,
	token.EQUAL_NOT:     COMPARE,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
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
	p.registerPrefixFunction(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixFunction(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFunction(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFunction(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFunction(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefixFunction(token.LBRACE, p.parseHashMap)
	p.registerPrefixFunction(token.STRING, p.parseStringLiteral)
	p.registerPrefixFunction(token.LPAREN, p.parseGrouping)

	// infix function
	p.registerInfixFunction(token.PLUS, p.parseInfixExpression)
	p.registerInfixFunction(token.MINUS, p.parseInfixExpression)
	p.registerInfixFunction(token.STAR, p.parseInfixExpression)
	p.registerInfixFunction(token.SLASH, p.parseInfixExpression)
	p.registerInfixFunction(token.GREATER, p.parseInfixExpression)
	p.registerInfixFunction(token.GREATER_EQUAL, p.parseInfixExpression)
	p.registerInfixFunction(token.LESSER, p.parseInfixExpression)
	p.registerInfixFunction(token.LESSER_EQUAL, p.parseInfixExpression)
	p.registerInfixFunction(token.EQUAL, p.parseInfixExpression)
	p.registerInfixFunction(token.EQUAL_NOT, p.parseInfixExpression)
	p.registerInfixFunction(token.LPAREN, p.parseFunctionCall)
	p.registerInfixFunction(token.LBRACKET, p.parseArrayIndex)

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

	case token.WHILE:
		return p.parseWhileStatement()

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
	// (1 + 2) * 3

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

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := &ast.PrefixExpression{Token: p.CurrentToken, Operator: p.CurrentToken.Literal}

	p.NextToken() // advance to the expression

	prefix.Right = p.parseExpression(PREFIX) // need to used other than LOWEST cause we are not the start of expression

	return prefix
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infix := &ast.InfixExpression{Token: p.CurrentToken, Left: left, Operator: p.CurrentToken.Literal}

	prec := precedence[p.CurrentToken.Type]
	p.NextToken() // advance to the right operand of the operator

	infix.Right = p.parseExpression(prec) // need to used other than LOWEST cause we are not the start of expression
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

func (p *Parser) parseWhileStatement() ast.Statement {
	whileStmt := &ast.WhileStatement{Token: p.CurrentToken}

	p.NextToken() // advance to `while` condition

	whileStmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	whileStmt.Body = p.parseBlockStatement().(*ast.BlockStatement)

	return whileStmt
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fnLit := &ast.FunctionLiteralExpression{Token: p.CurrentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.NextToken() // advance to function params

	fnLit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fnLit.Body = p.parseBlockStatement().(*ast.BlockStatement)

	return fnLit
}

func (p *Parser) parseFunctionParameters() []ast.Identifier {
	params := []ast.Identifier{}

	if p.currentTokenIs(token.RPAREN) {
		return params
	}

	param := ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.NextToken() // consume the `,` token
		p.NextToken() // advance to next expression in the list of expression

		param = ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
		params = append(params, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseFunctionCall(left ast.Expression) ast.Expression {
	fnCall := &ast.FunctionCallExpression{Token: p.CurrentToken, Function: left}

	p.NextToken() // advance to args list

	fnCall.Args = p.parseExpressionList(token.RPAREN)

	return fnCall
}

func (p *Parser) parseExpressionList(end token.TokenType) *ast.ExpressionList {
	exprList := &ast.ExpressionList{Token: p.CurrentToken}
	exprList.List = []ast.Expression{}

	if p.currentTokenIs(end) {
		return exprList
	}

	args := p.parseExpression(LOWEST)
	exprList.List = append(exprList.List, args)

	for p.peekTokenIs(token.COMMA) {
		p.NextToken() // consume the `,` token
		p.NextToken() // advance to next expression in the list of expression

		args := p.parseExpression(LOWEST)
		exprList.List = append(exprList.List, args)
	}

	if !p.expectPeek(end) {
		return nil
	}

	return exprList
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arrExpr := &ast.ArrayLiteralExpression{Token: p.CurrentToken}
	p.NextToken() // advance to array elements

	arrExpr.Elements = p.parseExpressionList(token.RBRACKET)

	return arrExpr
}

func (p *Parser) parseArrayIndex(left ast.Expression) ast.Expression {
	arrIndex := &ast.ArrayIndexExpression{Token: p.CurrentToken, Array: left}
	p.NextToken() // advance to index expression
	arrIndex.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return arrIndex
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteralExpression{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) parseHashMap() ast.Expression {
	hashMap := &ast.HashmapLiteralExpression{Token: p.CurrentToken}

	p.NextToken() // advance to expression

	hashMap.Map = p.parseHashmapExpressionList()
	return hashMap
}

func (p *Parser) parseHashmapExpressionList() map[ast.Expression]ast.Expression {
	hashMap := map[ast.Expression]ast.Expression{}

	if p.currentTokenIs(token.RBRACE) {
		return hashMap
	}

	key := p.parseExpression(LOWEST)

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.NextToken() // advance to next expression
	val := p.parseExpression(LOWEST)

	hashMap[key] = val

	for p.peekTokenIs(token.COMMA) {
		p.NextToken() // consume the `,`
		p.NextToken() // advance to next expression

		key = p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.NextToken() //advance to next expression

		val = p.parseExpression(LOWEST)
		hashMap[key] = val
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hashMap
}

func (p *Parser) parseGrouping() ast.Expression {

	// !(1 + 2)
	// (1 + 2) * 3
	// (1 + 2 * 3)
	p.NextToken() // advance to the expression
	left := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return left
}
