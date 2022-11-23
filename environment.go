package main

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment() Environment {
	return Environment{
		values:    make(map[string]any),
		enclosing: nil,
	}
}

/*
	Define a variable in an environment, used with assigning later on
*/
func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

/*
	Get a variable from an environment
*/
func (e *Environment) get(name Token) (any, error) {
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
func (e *Environment) assign(name Token, value any) error {
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
