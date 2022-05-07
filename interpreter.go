package main

import (
	"errors"
	"fmt"
	"reflect"
)

type interpreter struct{}

func (i *interpreter) interpret(expression Expr, l *Lox) {
	value, err := i.evaluate(expression)
	if err != nil {
		l.runtimeError(err)
	}
	fmt.Println(i.stringify(value))
}

func (i *interpreter) evaluate(expr Expr) (interface{}, error) {
	switch e := expr.(type) {
	case Literal:
		return e.value, nil
	case Unary:
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
	case Binary:
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
	}
	return nil, ParseError{errors.New("unreachable code error")}
}

func (i *interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if _, ok := obj.(bool); ok {
		return obj.(bool)
	}
	return true
}

func (i *interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *interpreter) checkNumberOperand(operator Token, operand interface{}) (float64, error) {
	o, ok := operand.(float64)
	if !ok {
		return 0, &RuntimeError{operator, "operand must be a number."}
	}
	return o, nil
}

func (i *interpreter) checkNumberOperands(operator Token, left, right interface{}) (float64, float64, error) {
	l, ok := left.(float64)
	r, ok2 := right.(float64)
	if !ok || !ok2 {
		return 0, 0, &RuntimeError{operator, "operands must be a number."}
	}
	return l, r, nil
}

func (i *interpreter) stringify(object interface{}) string {
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
