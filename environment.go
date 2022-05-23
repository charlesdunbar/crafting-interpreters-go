package main

type Environment struct {
	values map[string]interface{}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (interface{}, error) {
	if val, ok := e.values[name.lexeme]; ok {
		return val, nil
	}
	return nil, RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}

