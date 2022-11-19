package main

type LoxCallable interface {
	call(*interpreter, []any) any
	arity() int
}
