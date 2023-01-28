//TODO : parse hashmap is not working
package eval

import (
	"Klang/ast"
	"Klang/environment"
	"Klang/object"
	"fmt"
	"log"
)

var (
	NILL = &object.Nill{}
)

func Eval(node ast.Node, env *environment.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: float64(node.Value)}

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.StringLiteralExpression:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.Value}

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ArrayLiteralExpression:
		return evalArrayLiteralExpression(node, env)

	case *ast.HashmapLiteralExpression:
		return evalHashMapLiteralExpression(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	default:
		fmt.Println("HITTIN DEFAULT CASE")
		return NILL
	}
}

func evalProgram(statements []ast.Statement, env *environment.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)
	}

	return result
}

func evalInfixExpression(node *ast.InfixExpression, env *environment.Environment) object.Object {
	left := Eval(node.Left, env).(*object.Integer).Value
	right := Eval(node.Right, env).(*object.Integer).Value

	switch node.Operator {
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		return &object.Integer{Value: left / right}

	default:
		return &object.Integer{Value: 0}
	}
}

func evalLetStatement(node *ast.LetStatement, env *environment.Environment) object.Object {
	val := Eval(node.Value, env)
	env.Set(node.Name.Value, val)
	return NILL
}

func evalIdentifier(node *ast.Identifier, env *environment.Environment) object.Object {
	if val := env.Get(node.Value); val != nil {
		return val
	}

	return NILL
}

func evalArrayLiteralExpression(node *ast.ArrayLiteralExpression, env *environment.Environment) object.Object {
	objects := []object.Object{}

	for _, elem := range node.Elements.List {
		obj := Eval(elem, env)
		objects = append(objects, obj)
	}

	arr := &object.Array{Value: objects}
	return arr
}

func evalIndexExpression(node *ast.IndexExpression, env *environment.Environment) object.Object {
	ident := Eval(node.Ident, env)
	index := Eval(node.Index, env)

	switch ident.Type() {
	case object.OBJECT_ARRAY:
		array := ident.(*object.Array).Value
		idx := index.(*object.Integer).Value

		arrLen := len(array) - 1

		if idx < 0 || idx > int64(arrLen) {
			return NILL
		}

		return array[idx]

	case object.OBJECT_HASHMAP:
		hash := ident.(*object.HashMap).Value
		idx := index

		if val, ok := hash[idx]; ok {
			return val
		}

		return NILL

	default:
		return NILL
	}
}

func evalPrefixExpression(node *ast.PrefixExpression, env *environment.Environment) object.Object {
	switch node.Operator {
	case "!":
		log.Fatal("not implemented yet")
		return NILL

	case "-":
		val := Eval(node.Right, env).(*object.Integer).Value
		return &object.Integer{Value: -val}

	default:
		log.Fatalf("Unknown prefix operator %s\n", node.Operator)
		return NILL

	}
}

func evalHashMapLiteralExpression(node *ast.HashmapLiteralExpression, env *environment.Environment) object.Object {
	hash := make(map[object.Object]object.Object)

	for k, v := range node.Map {
		key := Eval(k, env)
		val := Eval(v, env)
		hash[key] = val
	}

	hashMap := &object.HashMap{Value: hash}
	return hashMap
}
