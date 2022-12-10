package main

type LoxClass struct {
	name    string
	methods map[string]LoxFunction
}

func NewLoxClass(name string, methods map[string]LoxFunction) LoxClass {
	return LoxClass{
		name:    name,
		methods: methods,
	}
}

func (l LoxClass) findMethod(name string) LoxFunction {
	if _, ok := l.methods[name]; ok {
		return l.methods[name]
	}
	return LoxFunction{}
}

func (l LoxClass) String() string {
	return l.name
}

func (l LoxClass) call(inter *interpreter, args []any) (any, error) {
	instance := NewLoxInstance(l)
	return instance, nil
}

func (l LoxClass) arity() int {
	return 0
}
