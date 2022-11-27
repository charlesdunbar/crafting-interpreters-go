package main

// import (
// 	"container/list"
// )

// type Resolver struct {
// 	interpreter interpreter
// 	scopes      list.List
// 	lox         *Lox
// }

// func (r Resolver) NewResolver() Resolver {
// 	return Resolver{
// 		// Is this the right thing I need?
// 		interpreter: *NewInterpreter(),
// 		scopes:      *list.New(),
// 	}
// }

// func (r Resolver) stmt_resolve(stmt Stmt) error {
// 	switch t := stmt.(type) {
// 	case *Block:
// 		r.beginScope()
// 		err := r.resolve_stmts(t.statements)
// 		if err != nil {
// 			return err
// 		}
// 		r.endScope()
// 	case *Var:
// 		r.declare(t.name)
// 		if t.initializer != nil {
// 			r.expr_resolve(t.initializer)
// 		}
// 		r.define(t.name)
// 	}
// 	return nil
// }

// func (r Resolver) expr_resolve(expr Expr) error {
// 	switch t := expr.(type) {
// 	case *Variable:
// 		front, ok := r.scopes.Front().Value.(map[string]bool)
// 		if !ok {
// 			panic("variable in expr_resolve can't cast correctly, scopes should only have map[string]bool types")
// 		}
// 		if r.scopes.Len() == 0 && !front[t.name.lexeme] {
// 			tokenError(t.name, "Can't read local variable in its own initializer.")
// 		}
// 		r.resolveLocal(t, t.name)
// 	}
// 	return nil
// }

// func (r Resolver) resolve_stmts(stmts []Stmt) error {
// 	for _, stmt := range stmts {
// 		err := r.stmt_resolve(stmt)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // Stack fun from https://medium.com/@dinesht.bits/stack-queue-implementations-in-golang-1136345036b4
// func (r Resolver) beginScope() {
// 	r.scopes.PushBack(make(map[string]bool))
// }

// func (r Resolver) endScope() {
// 	r.scopes.Remove(r.scopes.Back())
// }

// func (r Resolver) declare(name Token) {
// 	if r.scopes.Len() == 0 {
// 		return
// 	}
// 	scope, ok := r.scopes.Front().Value.(map[string]bool)
// 	if !ok {
// 		panic("declare somehow can't cast correctly, scopes should only have map[string]bool types")
// 	}
// 	scope[name.lexeme] = false
// }

// func (r Resolver) define(name Token) {
// 	if r.scopes.Len() == 0 {
// 		return
// 	}
// 	scope, ok := r.scopes.Front().Value.(map[string]bool)
// 	if !ok {
// 		panic("define somehow can't cast correctly, scopes should only have map[string]bool types")
// 	}
// 	scope[name.lexeme] = true

// }

// func (r Resolver) resolveLocal(expr Expr, name Token) {
// 	for i := r.scopes.Len() - 1; i >= 0; i-- {
// 		scope, ok := r.scopes.Front().Value.(map[string]bool)
// 		if !ok {
// 			panic("resolveLocal somehow can't cast correctly, scopes should only have map[string]bool types")
// 		}
// 		if _, ok := scope[name.lexeme]; ok {
// 			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
// 			return
// 		}
// 	}
// }
