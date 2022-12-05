package main

type LoxClass struct {
	name string
}

func NewLoxClass(name string) LoxClass {
	return LoxClass{
		name: name,
	}
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
