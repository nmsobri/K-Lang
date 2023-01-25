package ast

import (
	"Klang/token"
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	Statement()
}

type Expression interface {
	Node
	Expression()
}

type Program struct {
	Statements []Statement
}

// -----------------------------
// Let Statement
// -----------------------------
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Token.Literal)
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())

	return out.String()
}

func (ls *LetStatement) Statement() {}

// -----------------------------
// Return Statement
// -----------------------------
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Literal)
	out.WriteString(" ")
	out.WriteString(rs.ReturnValue.String())

	return out.String()
}

func (rs *ReturnStatement) Statement() {}

// -----------------------------
// Expression Statement
// -----------------------------
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(es.Expression.String())
	out.WriteString(")")

	return out.String()
}

func (es *ExpressionStatement) Statement() {}

// -----------------------------
// Identifier Expression
// -----------------------------
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) Expression() {}

// -----------------------------
// Integer Literal Expression
// -----------------------------
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

func (il *IntegerLiteral) Expression() {}

// -----------------------------
// Float Literal Expression
// -----------------------------
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%.2f", fl.Value)
}

func (fl *FloatLiteral) Expression() {}

// -----------------------------
// Boolean Literal Expression
// -----------------------------
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(fmt.Sprintf("%t", bl.Value))
	out.WriteString(")")

	return out.String()
}

func (bl *BooleanLiteral) Expression() {}

// -----------------------------
// Infix Expression
// -----------------------------
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) Expression() {}

// -----------------------------
// If Expression
// -----------------------------
type IfExpression struct {
	Token     token.Token
	Condition Expression
	IfArm     *BlockStatement
	ElseArm   *BlockStatement
}

func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ife.Token.Literal)
	out.WriteString(ife.Condition.String())
	out.WriteString(ife.IfArm.String())

	if ife.ElseArm != nil {
		out.WriteString("else")
		out.WriteString(ife.ElseArm.String())
	}

	return out.String()
}

func (ife *IfExpression) Expression() {}

// -----------------------------
// Block Statement
// -----------------------------
type BlockStatement struct {
	Token      token.Token // the `{`
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
		out.WriteString(";")
	}

	out.WriteString("}")

	return out.String()
}

func (bs *BlockStatement) Expression() {}

// -----------------------------
// While Statement
// -----------------------------
type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}

func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.Token.Literal)
	out.WriteString(ws.Condition.String())
	out.WriteString("{")
	out.WriteString(ws.Body.String())
	out.WriteString("}")

	return out.String()
}

func (ws *WhileStatement) Statement() {}

// -----------------------------
// Function Literal Expression
// -----------------------------
type FunctionExpression struct {
	Token      token.Token
	Parameters []Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) TokenLiteral() string {
	return fe.Token.Literal
}

func (fe *FunctionExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fe.Token.Literal)
	out.WriteString("(")

	params := []string{}

	for _, param := range fe.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(strings.Join(params, ", "))

	out.WriteString(")")
	out.WriteString(fe.Body.String())

	return out.String()
}

func (fe *FunctionExpression) Expression() {}

// -----------------------------
// Expression List
// -----------------------------
type ExpressionList struct {
	Token token.Token // the `(` token
	List  []Expression
}

func (el *ExpressionList) TokenLiteral() string {
	return el.Token.Literal
}

func (el *ExpressionList) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	args := []string{}
	for _, stmt := range el.List {
		args = append(args, stmt.String())
	}

	out.WriteString(strings.Join(args, ", "))

	out.WriteString(")")

	return out.String()
}

func (el *ExpressionList) Expression() {}

// -----------------------------
// Function Call Expression
// -----------------------------
type FunctionCallExpression struct {
	Token    token.Token
	Function Expression // the function itself
	Args     *ExpressionList
}

func (fc *FunctionCallExpression) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc *FunctionCallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fc.Function.String())
	out.WriteString("(")

	args := []string{}
	for _, expr := range fc.Args.List {
		args = append(args, expr.String())
	}

	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (fc *FunctionCallExpression) Expression() {}
