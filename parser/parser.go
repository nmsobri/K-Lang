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
)

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

	p.registerPrefixFunction(token.INTEGER, p.parseInteger)

	// prime the tokens
	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) NextToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.Lexer.NextToken()
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
	p.NextToken() // advance to identifier

	letStmt.Name = p.parseIdentifier().(*ast.Identifier)

	p.NextToken() // advance to value expression
	letStmt.Value = p.parseExpression(LOWEST)

	return letStmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	returnStmt := &ast.ReturnStatement{Token: p.CurrentToken}
	p.NextToken() // advance to value expression
	returnStmt.ReturnValue = p.parseExpression(LOWEST)
	return returnStmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expressionStmt := &ast.ExpressionStatement{Token: p.CurrentToken}
	expressionStmt.Expression = p.parseExpression(LOWEST)
	return expressionStmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	p.NextToken()
	return ident
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.getPrefixFunction()

	if prefix == nil {
		return nil
	}

	left := prefix()
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

func (p *Parser) getPrefixFunction() prefixFunc {
	if fn, ok := p.prefixFunc[p.CurrentToken.Type]; ok {
		return fn
	}

	return nil
}

func (p *Parser) getInfixFunction() infixFunc {
	if fn, ok := p.infixFunc[p.CurrentToken.Type]; ok {
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
