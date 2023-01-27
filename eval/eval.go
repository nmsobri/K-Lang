package eval

import (
	"Klang/ast"
)

func Eval(node ast.Node) int {
	switch node := node.(type) {
	case *ast.Program:
		evaluated := evalProgram(node.Statements)
		return evaluated

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return int(node.Value)

	case *ast.InfixExpression:
		return evalInfixExpression(node)

	default:
		return 99
	}
}

func evalProgram(statements []ast.Statement) int {
	var result int

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func evalInfixExpression(node *ast.InfixExpression) int {
	left := Eval(node.Left)
	right := Eval(node.Right)

	switch node.Operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right

	default:
		return 0
	}
}
