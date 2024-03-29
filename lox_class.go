package main

import "errors"

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]LoxFunction) LoxClass {
	return LoxClass{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (l LoxClass) findMethod(name string) (LoxFunction, error) {
	if _, ok := l.methods[name]; ok {
		return l.methods[name], nil
	}

	if l.superclass != nil {
		return l.superclass.findMethod(name)
	}

	return LoxFunction{}, errors.New("method not found")
}

func (l LoxClass) String() string {
	return l.name
}

func (l LoxClass) call(inter *interpreter, args []any) (any, error) {
	instance := NewLoxInstance(l)
	initalizer, err := l.findMethod("init")
	if err == nil {
		initalizer.bind(instance).call(inter, args)
	}
	return instance, nil
}

func (l LoxClass) arity() int {
	initalizer, err := l.findMethod("init")
	if err != nil {
		return 0
	}
	return initalizer.arity()
}
