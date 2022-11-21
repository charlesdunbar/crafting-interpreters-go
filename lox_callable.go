package main

type LoxCallable interface {
	call(*interpreter, []any) (any, error)
	arity() int
}
