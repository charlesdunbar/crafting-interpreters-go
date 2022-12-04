package main

import (
	"container/list"
)

type FunctionType int64

const (
	NONE FunctionType = iota
	FUNCTION
)

type Resolver struct {
	interpreter     interpreter
	scopes          list.List
	currentFunction FunctionType
}

func (r Resolver) NewResolver() Resolver {
	return Resolver{
		// Is this the right thing I need?
		interpreter:     *NewInterpreter(),
		scopes:          *list.New(),
		currentFunction: NONE,
	}
}

func (r Resolver) stmt_resolve(stmt Stmt) error {
	switch t := stmt.(type) {
	case *Block:
		r.beginScope()
		err := r.resolve_stmts(t.statements)
		if err != nil {
			return err
		}
		r.endScope()
	case *Expression:
		r.expr_resolve(t.expression)
	case *Function:
		r.declare(t.name)
		r.define(t.name)
		r.resolveFunction(*t, FUNCTION)
	case *If:
		r.expr_resolve(t.condition)
		r.stmt_resolve(t.thenBranch)
		if t.elseBranch != nil {
			r.stmt_resolve(t.elseBranch)
		}
	case *Print:
		r.expr_resolve(t.expression)
	case *Return:
		if r.currentFunction == NONE {
			tokenError(t.keyword, "Can't return from top-level code.")
		}
		if t.value != nil {
			r.expr_resolve(t.value)
		}
	case *Var:
		r.declare(t.name)
		if t.initializer != nil {
			r.expr_resolve(t.initializer)
		}
		r.define(t.name)
	case *While:
		r.expr_resolve(t.condition)
		r.stmt_resolve(t.body)
	}
	return nil
}

func (r Resolver) expr_resolve(expr Expr) error {
	switch t := expr.(type) {
	case *Assign:
		err := r.expr_resolve(t.value)
		if err != nil {
			return err
		}
		r.resolveLocal(t, t.name)
	case *Binary:
		r.expr_resolve(t.left)
		r.expr_resolve(t.right)
	case *Call:
		r.expr_resolve(t.callee)
		for _, arg := range t.arguments {
			r.expr_resolve(arg)
		}
	case *Grouping:
		r.expr_resolve(t.expression)
	case *Literal:
		return nil
	case *Logical:
		r.expr_resolve(t.left)
		r.expr_resolve(t.right)
	case *Unary:
		r.expr_resolve(t.right)
	case *Variable:
		if r.scopes.Len() != 0 {
			front, ok := r.scopes.Front().Value.(map[string]bool)
			if !ok {
				panic("variable in expr_resolve can't cast correctly, scopes should only have map[string]bool types")
			}
			if !front[t.name.lexeme] {
				tokenError(t.name, "Can't read local variable in its own initializer.")
			}
		}
		r.resolveLocal(t, t.name)
	}
	return nil
}

func (r Resolver) resolve_stmts(stmts []Stmt) error {
	for _, stmt := range stmts {
		err := r.stmt_resolve(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(function Function, t FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = t
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolve_stmts(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

// Stack fun from https://medium.com/@dinesht.bits/stack-queue-implementations-in-golang-1136345036b4
func (r *Resolver) beginScope() {
	r.scopes.PushBack(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Remove(r.scopes.Back())
}

func (r Resolver) declare(name Token) {
	if r.scopes.Len() == 0 {
		return
	}
	scope, ok := r.scopes.Front().Value.(map[string]bool)
	if !ok {
		panic("declare somehow can't cast correctly, scopes should only have map[string]bool types")
	}
	if _, ok := scope[name.lexeme]; ok {
		tokenError(name, "Already a variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

func (r Resolver) define(name Token) {
	if r.scopes.Len() == 0 {
		return
	}
	scope, ok := r.scopes.Front().Value.(map[string]bool)
	if !ok {
		panic("define somehow can't cast correctly, scopes should only have map[string]bool types")
	}
	scope[name.lexeme] = true

}

func (r Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		scope, ok := r.scopes.Front().Value.(map[string]bool)
		if !ok {
			panic("resolveLocal somehow can't cast correctly, scopes should only have map[string]bool types")
		}
		if _, ok := scope[name.lexeme]; ok {
			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
			return
		}
	}
}
