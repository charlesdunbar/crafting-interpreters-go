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
	env := &Environment{
		values: make(map[string]any),
		// This is nil and probably shouldn't be
		enclosing: inter.globals,
	}
	for i := 0; i < len(l.declaration.params); i++ {
		// Try to not store an Expr in the environment, but what the Expr evalutes to
		// var err error
		// if a, ok := args[i].(Expr); ok {
		// 	switch e := a.(type) {
		// 	case *Literal:
		// 		fmt.Printf("Evaluating %+v before storing it in an environment\n", e)
		// 		args[i], err = inter.evaluate(e)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 	case *Variable:
		// 		fmt.Printf("Evaluating %+v before storing it in an environment\n", e)
		// 		args[i], err = inter.evaluate(e)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 	}
		// }
		env.define(
			l.declaration.params[i].lexeme,
			// Stores Statements or Expressions, such as Literal or Function
			// This causes problems as it doesn't mirror behavior as a variable does (which stores the evaluate value of an initalizer)
			args[i],
		)
	}
	// https://stackoverflow.com/a/44543748
	err := (&interpreter{}).executeBlock(l.declaration.body, env)
	fmt.Printf("About to return an error maybe with %+v\n", err)
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
