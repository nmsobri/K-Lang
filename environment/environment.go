package environment

import "Klang/object"

type Environment struct {
	Vars   map[string]object.Object
	Parent *Environment
}

func New() *Environment {
	env := Environment{}
	env.Vars = make(map[string]object.Object)
	env.Parent = nil

	return &env
}

func NewWithParent(parent *Environment) *Environment {
	env := New()
	env.Parent = parent
	return env
}

func (env *Environment) Set(key string, val object.Object) object.Object {
	env.Vars[key] = val
	return val
}

func (env *Environment) Get(key string) object.Object {
	if val, ok := env.Vars[key]; ok {
		return val
	}

	if env.Parent != nil {
		if val := env.Parent.Get(key); val != nil {
			return val
		}
	}

	return nil
}
