package environment

import "Klang/object"

type Environment struct {
	Vars map[string]object.Object
}

func New() *Environment {
	env := Environment{}
	env.Vars = make(map[string]object.Object)

	return &env
}

func (env *Environment) Set(key string, val object.Object) object.Object {
	env.Vars[key] = val
	return val
}

func (env *Environment) Get(key string) object.Object {
	if val, ok := env.Vars[key]; ok {
		return val
	}

	return nil
}
