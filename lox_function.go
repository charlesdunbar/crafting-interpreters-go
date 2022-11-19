package main

import "fmt"

type LoxFunction struct {
	declaration Function
}

func NewLoxFunction(dec Function) *LoxFunction {
	return &LoxFunction{
		dec,
	}
}

// Implement LoxCallable
func (l LoxFunction) call(inter *interpreter, args []any) any {
	env := &Environment{
		values:    make(map[string]any),
		enclosing: inter.globals,
	}
	for i := 0; i < len(l.declaration.params); i++ {
		env.define(
			l.declaration.params[i].lexeme,
			args[i],
		)
	}
	// https://stackoverflow.com/a/44543748
	(&interpreter{}).executeBlock(l.declaration.body, env)
	return nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.declaration.name)
}
