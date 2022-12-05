package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type interpreter struct {
	globals     Environment
	environment *Environment
	locals      map[Expr]int
}

// Built-in clock functionality
type clock struct{}

func (c clock) arity() int {
	return 0
}

func (c clock) call(int *interpreter, args []any) (any, error) {
	return float64(time.Now().UnixMilli() / 1000), nil
}

func (c clock) String() string {
	return "<native fn>"
}

func NewInterpreter() *interpreter {
	global := NewEnvironment()
	env := &global

	global.define("clock", clock{})
	return &interpreter{
		globals:     global,
		environment: env,
		locals:      make(map[Expr]int),
	}
}

func (i interpreter) interpret(statements []Stmt) error {
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
	case *Class:
		i.environment.define(t.name.lexeme, nil)
		c := NewLoxClass(t.name.lexeme)
		err := i.environment.assign(t.name, c)
		if err != nil {
			return err
		}
	case *Expression:
		_, err := i.evaluate(t.expression)
		if err != nil {
			return err
		}
	case *Function:
		function := LoxFunction{*t, *i.environment}
		i.environment.define(t.name.lexeme, function)
	case *If:
		cond, err := i.evaluate(t.condition)
		if err != nil {
			return err
		}
		if i.isTruthy(cond) {
			err := i.execute(t.thenBranch)
			if err != nil {
				return err
			}
		} else if t.elseBranch != nil {
			err := i.execute(t.elseBranch)
			if err != nil {
				return err
			}
		}
	case *Print:
		value, err := i.evaluate(t.expression)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", i.stringify(value))
	case *Return:
		var value any
		var err error
		if t.value != nil {
			value, err = i.evaluate(t.value)
			if err != nil {
				return err
			}
		}
		return NewReturnError(value)
	case *While:
		// Can't chain i.evalute(t.condition) by itself, so set it as initalizer and incrementer
		for val, err := i.evaluate(t.condition); i.isTruthy(val); val, err = i.evaluate(t.condition) {
			if err != nil {
				return err
			}
			err := i.execute(t.body)
			if err != nil {
				return err
			}
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

func (i *interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
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
		if distance, ok := i.locals[expr]; ok {
			i.environment.assignAt(distance, e.name, value)
		} else {
			err = i.globals.assign(e.name, value)
			if err != nil {
				return nil, err
			}
		}

		return value, nil
	case *Call:
		callee, err := i.evaluate(e.callee)
		if err != nil {
			return nil, err
		}
		var arguments []any
		for _, arg := range e.arguments {
			eval, err := i.evaluate(arg)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, eval)
		}

		// Don't allow trying to call "foobar"()
		local_func, ok := callee.(LoxCallable)
		if !ok {
			return nil, NewRuntimeError(e.paren, fmt.Sprintf("Can only call functions and classes. Got data of type %T", callee))
		}

		// Check arity of the function
		if len(arguments) != local_func.arity() {
			return nil, NewRuntimeError(e.paren, fmt.Sprintf("Expected %d arguments but got %d.\n", local_func.arity(), len(arguments)))
		}
		ret, err := local_func.call(i, arguments)
		if err != nil {
			return nil, err
		}
		return ret, nil
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
			if l, ok := left.(float64); ok {
				if r, ok := right.(float64); ok {
					return l + r, nil
				}
			} else if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r, nil
				}
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
		return i.lookUpVariable(e.name, e)
	}
	return nil, ParseError{errors.New("unreachable code error")}
}

func (i *interpreter) lookUpVariable(name Token, expr Expr) (any, error) {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, name.lexeme), nil
	}
	return i.globals.get(name)

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
