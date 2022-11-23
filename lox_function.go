package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration Function
	closure Environment
}

func NewLoxFunction(dec Function, clo Environment) *LoxFunction {
	return &LoxFunction{
		declaration: dec,
		closure: clo,
	}
}

// Implement LoxCallable
func (l LoxFunction) call(inter *interpreter, args []any) (any, error) {
	env := Environment{
		values: make(map[string]any),
		// This is nil and probably shouldn't be
		enclosing: &l.closure,
	}
	for i := 0; i < len(l.declaration.params); i++ {
		env.define(
			l.declaration.params[i].lexeme,
			args[i],
		)
	}
	err := inter.executeBlock(l.declaration.body, &env)
	if err != nil {
		switch e := err.(type) {
		// We actually want to use this as an exception to break early, not as an error that needs reporting
		case *ReturnError:
			return e.value, nil
		case *RuntimeError:
			return nil, e
		}
	}
	return nil, nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.declaration.name.lexeme)
}
