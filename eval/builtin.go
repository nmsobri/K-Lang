package eval

import (
	"Klang/object"
	"fmt"
	"strings"
)

type BuiltinFn func(args ...object.Object) object.Object

func (bf BuiltinFn) Inspect() string {
	return ""
}

func (bf BuiltinFn) Type() object.ObjectType {
	return object.OBJECT_BUILTIN
}

var builtins = map[string]BuiltinFn{
	"len": func(args ...object.Object) object.Object {
		length := len(args)

		if length != 1 {
			return NILL
		}

		if args[0].Type() != object.OBJECT_ARRAY {
			return NILL
		}

		arr := args[0].(*object.Array)
		arrLen := len(arr.Value)

		return &object.Integer{Value: int64(arrLen)}
	},
	"print": func(args ...object.Object) object.Object {
		arguments := []string{}

		for _, arg := range args {
			arguments = append(arguments, arg.Inspect())
		}

		fmt.Printf("%s\n", strings.Join(arguments, " "))
		return NILL
	},
}
