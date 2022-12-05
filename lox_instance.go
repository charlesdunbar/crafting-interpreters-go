package main

import "fmt"

type LoxInstance struct {
	klass LoxClass
}

func NewLoxInstance(klass LoxClass) LoxInstance {
	return LoxInstance{klass: klass}
}

func (l LoxInstance) String() string {
	return fmt.Sprintf("%s instance", l.klass.name)
}
