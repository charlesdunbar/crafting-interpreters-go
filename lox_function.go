package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration   Function
	closure       Environment
	isInitializer bool
}

func NewLoxFunction(dec Function, clo Environment, init bool) LoxFunction {
	return LoxFunction{
		declaration: dec,
		closure:     clo,
		isInitializer: init,
	}
}

func (l LoxFunction) bind(instance LoxInstance) LoxFunction {
	environment := Environment{
		values:    make(map[string]any),
		enclosing: &l.closure,
	}
	environment.define("this", instance)
	return NewLoxFunction(l.declaration, environment, l.isInitializer)

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
			if l.isInitializer {
				return l.closure.getAt(0, "this"), nil
			}
			return e.value, nil
		case *RuntimeError:
			return nil, e
		}
	}
	if l.isInitializer {
		return l.closure.getAt(0, "this"), nil
	}
	return nil, nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.declaration.name.lexeme)
}
