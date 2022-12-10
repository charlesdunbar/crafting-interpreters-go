package main

import (
	"fmt"
	"reflect"
)

type LoxInstance struct {
	klass  LoxClass
	fields map[string]any
}

func NewLoxInstance(klass LoxClass) LoxInstance {
	return LoxInstance{
		klass:  klass,
		fields: make(map[string]any),
	}
}

func (l *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.klass.name)
}

func (l *LoxInstance) get(name Token) (any, error) {
	if _, ok := l.fields[name.lexeme]; ok {
		return l.fields[name.lexeme], nil
	}

	method := l.klass.findMethod(name.lexeme)
	// If not an empty LoxFunction, return the method
	if !reflect.ValueOf(method).IsZero() {
		return method.bind(*l), nil
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined property %s.", name.lexeme))
}

func (l *LoxInstance) set(name Token, value any) {
	l.fields[name.lexeme] = value
}
