package main

import "golang.org/x/tools/go/analysis/passes/nilfunc"

type Environment struct {
	values map[string]interface{}
	enclosing *Environment
}

/*
	Define a variable in an environment, used with assigning later on
*/
func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

/*
	Get a variable from an environment
*/
func (e *Environment) get(name Token) (interface{}, error) {
	if val, ok := e.values[name.lexeme]; ok {
		return val, nil
	}

	/*
		Grab variable from an upper environment if it doesn't exist in this one
	*/
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}

/*
	Assign a variable to an environment, erroring if it's undefined.
*/
func (e *Environment) assign(name Token, value interface{}) error {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	/*
		Shadow a variable if one is defined in an upper environment
	*/
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return nil
	}

	/*
		Error if variable not defined anywhere
	*/
	return RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}
