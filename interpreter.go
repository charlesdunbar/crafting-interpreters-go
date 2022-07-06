package main

import (
	"errors"
	"fmt"
	"reflect"
)

type interpreter struct {
	environment *Environment
}

func (i *interpreter) interpret(statements []Stmt) error {
	for _, s := range statements {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Visit statement replacement
func (i *interpreter) execute(stmt Stmt) error {
	switch t := stmt.(type) {
	case *Block:
		err := i.executeBlock(t.statements, &Environment{values: make(map[string]any), enclosing: i.environment})
		if err != nil {
			return err
		}
	case *Expression:
		_, err := i.evaluate(t.expression)
		if err != nil {
			return err
		}
	case *If:
		cond, err := i.evaluate(t.condition)
		if err != nil {
			return err
		}
		if i.isTruthy(cond) {
			i.execute(t.thenBranch)
		} else if t.elseBranch != nil {
			i.execute(t.elseBranch)
		}
	case *Print:
		value, err := i.evaluate(t.expression)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", i.stringify(value))
	case *While:
		// Can't chain i.evalute(t.condition) by itself, so set it as initalizer and incrementer
		for val, err := i.evaluate(t.condition); i.isTruthy(val); val, err = i.evaluate(t.condition) {
			if err != nil {
				return err
			}
			i.execute(t.body)
		}
	case *Var:
		var value any
		var err error
		if t.initializer != nil {
			value, err = i.evaluate(t.initializer)
			if err != nil {
				return err
			}
		} else {
			value = nil
		}
		i.environment.define(t.name.lexeme, value)
	}
	return nil
}

func (i *interpreter) executeBlock(statements []Stmt, env *Environment) error {
	previous := i.environment
	// Mimic "finally" block
	defer func() { i.environment = previous }()

	i.environment = env
	for _, stmt := range statements {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Visit Expression replacement
func (i *interpreter) evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {
	case *Assign:
		value, err := i.evaluate(e.value)
		if err != nil {
			return nil, err
		}
		err = i.environment.assign(e.name, value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case *Literal:
		return e.value, nil
	case *Logical:
		left, err := i.evaluate(e.left)
		if err != nil {
			return nil, err
		}

		if e.operator.l_type == OR {
			// Short circuit true for OR
			if i.isTruthy(left) {
				return left, nil
			} else {
				// Short circuit false for AND
				if !i.isTruthy(left) {
					return left, nil
				}
			}
		}
		return i.evaluate(e.right)
	case *Unary:
		right, err := i.evaluate(e.right)
		if err != nil {
			return nil, ParseError{err}
		}

		switch e.operator.l_type {
		case BANG:
			return !i.isTruthy(right), nil
		case MINUS:
			r, err := i.checkNumberOperand(e.operator, right)
			if err != nil {
				return nil, err
			}
			return -r, nil
		}
		// Unreachable
		return nil, ParseError{err}
	case *Binary:
		left, err := i.evaluate(e.left)
		if err != nil {
			return nil, ParseError{err}
		}
		right, err := i.evaluate(e.right)
		if err != nil {
			return nil, ParseError{err}
		}

		switch e.operator.l_type {
		case GREATER:
			l, r, err := i.checkNumberOperands(e.operator, left, right)
			if err != nil {
				return nil, err
			}
			return l > r, nil
		case GREATER_EQUAL:
			l, r, err := i.checkNumberOperands(e.operator, left, right)
			if err != nil {
				return nil, err
			}
			return l >= r, nil
		case LESS:
			l, r, err := i.checkNumberOperands(e.operator, left, right)
			if err != nil {
				return nil, err
			}
			return l < r, nil
		case LESS_EQUAL:
			l, r, err := i.checkNumberOperands(e.operator, left, right)
			if err != nil {
				return nil, err
			}
			return l <= r, nil
		case BANG_EQUAL:
			return !i.isEqual(left, right), nil
		case EQUAL_EQUAL:
			return i.isEqual(left, right), nil
		case MINUS:
			return left.(float64) - right.(float64), nil
		case PLUS:
			if reflect.TypeOf(left).Kind().String() == "float64" &&
				reflect.TypeOf(right).Kind().String() == "float64" {
				return left.(float64) + right.(float64), nil
			}
			if reflect.TypeOf(left).Kind().String() == "string" &&
				reflect.TypeOf(right).Kind().String() == "string" {
				return left.(string) + right.(string), nil
			}
			return nil, &RuntimeError{e.operator, "operands must be two numbers or two strings."}
		case SLASH:
			return left.(float64) / right.(float64), nil
		case STAR:
			return left.(float64) * right.(float64), nil
		}
		// Unreachable
		return nil, ParseError{err}
	case *Variable:
		return i.environment.get(e.name)
	}
	return nil, ParseError{errors.New("unreachable code error")}
}

func (i *interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if _, ok := obj.(bool); ok {
		return obj.(bool)
	}
	return true
}

func (i *interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *interpreter) checkNumberOperand(operator Token, operand any) (float64, error) {
	o, ok := operand.(float64)
	if !ok {
		return 0, &RuntimeError{operator, "operand must be a number."}
	}
	return o, nil
}

func (i *interpreter) checkNumberOperands(operator Token, left, right any) (float64, float64, error) {
	l, ok := left.(float64)
	r, ok2 := right.(float64)
	if !ok || !ok2 {
		return 0, 0, &RuntimeError{operator, "operands must be a number."}
	}
	return l, r, nil
}

func (i *interpreter) stringify(object any) string {
	if object == nil {
		return "nil"
	}

	if reflect.TypeOf(object).Kind().String() == "float64" {
		text := fmt.Sprintf("%.1f", object)
		if text[len(text)-2:] == ".0" {
			text = text[:len(text)-2]
		}
		return text
	}
	return fmt.Sprintf("%v", object)

}
