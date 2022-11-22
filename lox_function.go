package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration Function
}

func NewLoxFunction(dec Function) *LoxFunction {
	return &LoxFunction{
		dec,
	}
}

// Implement LoxCallable
func (l LoxFunction) call(inter *interpreter, args []any) (any, error) {
	env := Environment{
		values: make(map[string]any),
		// This is nil and probably shouldn't be
		enclosing: &inter.globals,
	}
	for i := 0; i < len(l.declaration.params); i++ {
		env.define(
			l.declaration.params[i].lexeme,
			args[i],
		)
	}
	// https://stackoverflow.com/a/44543748
	err := inter.executeBlock(l.declaration.body, &env)
	//fmt.Printf("About to return an error maybe with %+v\n", err)
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
