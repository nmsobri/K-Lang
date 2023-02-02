package eval

import (
	"Klang/ast"
	"Klang/object"
	"log"
)

var (
	NILL  = &object.Nill{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
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
		if node.Value {
			return TRUE
		}
		return FALSE

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

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)

	case *ast.FunctionLiteralExpression:
		return evalFunctionLiteralExpression(node, env)

	case *ast.FunctionCallExpression:
		return evalFunctionCallExpression(node, env)

	case *ast.ExpressionList:
		return evalExpressionList(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	default:
		log.Fatalf("Unhandled case for: %T", node)
		return NILL
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)

		if result.Type() == object.OBJECT_RETURN {
			return result
		}
	}

	return result
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
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

	case ">":
		return &object.Boolean{Value: left > right}

	case ">=":
		return &object.Boolean{Value: left >= right}

	case "<":
		return &object.Boolean{Value: left < right}

	case "<=":
		return &object.Boolean{Value: left <= right}

	case "==":
		return &object.Boolean{Value: left == right}

	case "!=":
		return &object.Boolean{Value: left != right}

	default:
		log.Fatalf("Unknown infix operator %s\n", node.Operator)
		return NILL
	}
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	env.Set(node.Name.Value, val)
	return NILL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val := env.Get(node.Value); val != nil {
		return val
	}

	log.Printf("undefined identifier: %s\n", node.Value)
	return NILL
}

func evalArrayLiteralExpression(node *ast.ArrayLiteralExpression, env *object.Environment) object.Object {
	objects := []object.Object{}

	for _, elem := range node.Elements.List {
		obj := Eval(elem, env)
		objects = append(objects, obj)
	}

	arr := &object.Array{Value: objects}
	return arr
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
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
		idx := index.(object.Hashable)

		if val, ok := hash[idx.Hashkey()]; ok {
			return val
		}

		return NILL

	default:
		return NILL
	}
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	switch node.Operator {
	case "!":
		val := Eval(node.Right, env)
		boolean := isTruthy(val)
		return &object.Boolean{Value: !boolean}

	case "-":
		val := Eval(node.Right, env).(*object.Integer).Value
		return &object.Integer{Value: -val}

	default:
		log.Fatalf("Unknown prefix operator %s\n", node.Operator)
		return NILL

	}
}

func evalHashMapLiteralExpression(node *ast.HashmapLiteralExpression, env *object.Environment) object.Object {
	hashMap := make(map[object.Hash]object.Object)

	for k, v := range node.Map {
		key := Eval(k, env)
		val := Eval(v, env)

		hash, ok := key.(object.Hashable)

		if !ok {
			log.Fatalf("invalid key type: %T", key)
			return NILL
		}

		hashMap[hash.Hashkey()] = val
	}

	return &object.HashMap{Value: hashMap}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)

	if isTruthy(condition) {
		return Eval(node.IfArm, env)
	}

	return Eval(node.ElseArm, env)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NILL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	return evalProgram(node.Statements, env)
}

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	var res object.Object

	for Eval(node.Condition, env).(*object.Boolean).Value {
		res = Eval(node.Body, env)
	}

	return res
}

func evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	env.Set(node.Ident.Value, val)
	return NILL
}

func evalExpressionList(node *ast.ExpressionList, env *object.Environment) object.Object {
	expressions := []object.Object{}

	for _, expr := range node.List {
		obj := Eval(expr, env)
		expressions = append(expressions, obj)
	}

	return &object.Array{Value: expressions}
}

func evalFunctionLiteralExpression(node *ast.FunctionLiteralExpression, env *object.Environment) object.Object {
	return &object.Function{Parameters: node.Parameters, Body: node.Body, Environment: env}
}

func evalFunctionCallExpression(node *ast.FunctionCallExpression, env *object.Environment) object.Object {
	fn := Eval(node.Function, env).(*object.Function)
	args := Eval(node.Args, env).(*object.Array)

	// start function own scope and inherit from outter scope
	fnEnv := object.NewEnvironmentWithParent(fn.Environment)

	// bind args to params
	for k, v := range args.Value {
		fnEnv.Set(fn.Parameters[k].Value, v)
	}

	return Eval(fn.Body, fnEnv)
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	value := Eval(node.ReturnValue, env)
	return &object.Return{Value: value}
}
