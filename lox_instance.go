package main

import (
	"fmt"
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

func (l LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.klass.name)
}

func (l *LoxInstance) get(name Token) (any, error) {
	if _, ok := l.fields[name.lexeme]; ok {
		return l.fields[name.lexeme], nil
	}

	method, err := l.klass.findMethod(name.lexeme)
	if err != nil {
		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined property %s.", name.lexeme))
	}
	return method.bind(*l), nil
}

func (l *LoxInstance) set(name Token, value any) {
	l.fields[name.lexeme] = value
}
