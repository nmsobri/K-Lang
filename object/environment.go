package object

type Environment struct {
	Vars   map[string]Object
	Parent *Environment
}

func NewEnvironment() *Environment {
	env := Environment{}
	env.Vars = make(map[string]Object)
	env.Parent = nil

	return &env
}

func NewEnvironmentWithParent(parent *Environment) *Environment {
	env := NewEnvironment()
	env.Parent = parent
	return env
}

func (env *Environment) Set(key string, val Object) Object {
	env.Vars[key] = val
	return val
}

func (env *Environment) Get(key string) Object {
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
