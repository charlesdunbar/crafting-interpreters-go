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
			// Stores Statements or Expressions, such as Literal or Function
			// This causes problems as it doesn't mirror behavior as a variable does (which stores the evaluate value of an initalizer)
			args[i],
		)
	}
	// https://stackoverflow.com/a/44543748
	err := (&interpreter{}).executeBlock(l.declaration.body, env)
	if err != nil {
		switch e := err.(type) {
		case ReturnError:
			return e.value
		}
	}
	return nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.declaration.name)
}
