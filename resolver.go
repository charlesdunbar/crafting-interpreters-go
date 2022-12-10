package main

type FunctionType int64

const (
	NONE FunctionType = iota
	FUNCTION
)

type Resolver struct {
	interpreter     interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func (r *Resolver) NewResolver() Resolver {
	return Resolver{
		interpreter:     *NewInterpreter(),
		scopes:          make([]map[string]bool, 0),
		currentFunction: NONE,
	}
}

func (r *Resolver) stmt_resolve(stmt Stmt) error {
	switch t := stmt.(type) {
	case *Block:
		r.beginScope()
		err := r.resolve_stmts(t.statements)
		if err != nil {
			return err
		}
		r.endScope()
	case *Class:
		r.declare(t.name)
		r.define(t.name)
	case *Expression:
		err := r.expr_resolve(t.expression)
		if err != nil {
			return nil
		}
	case *Function:
		r.declare(t.name)
		r.define(t.name)
		r.resolveFunction(*t, FUNCTION)
	case *If:
		err := r.expr_resolve(t.condition)
		if err != nil {
			return nil
		}

		err = r.stmt_resolve(t.thenBranch)
		if err != nil {
			return nil
		}
		if t.elseBranch != nil {
			err = r.stmt_resolve(t.elseBranch)
			if err != nil {
				return nil
			}
		}
	case *Print:
		err := r.expr_resolve(t.expression)
		if err != nil {
			return nil
		}
	case *Return:
		if r.currentFunction == NONE {
			tokenError(t.keyword, "Can't return from top-level code.")
		}
		if t.value != nil {
			err := r.expr_resolve(t.value)
			if err != nil {
				return nil
			}
		}
	case *Var:
		r.declare(t.name)
		if t.initializer != nil {
			err := r.expr_resolve(t.initializer)
			if err != nil {
				return nil
			}
		}
		r.define(t.name)
	case *While:
		err := r.expr_resolve(t.condition)
		if err != nil {
			return nil
		}

		err = r.stmt_resolve(t.body)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (r *Resolver) expr_resolve(expr Expr) error {
	switch t := expr.(type) {
	case *Assign:
		err := r.expr_resolve(t.value)
		if err != nil {
			return err
		}
		r.resolveLocal(t, t.name)
	case *Binary:
		err := r.expr_resolve(t.left)
		if err != nil {
			return nil
		}
		err = r.expr_resolve(t.right)
		if err != nil {
			return nil
		}
	case *Call:
		r.expr_resolve(t.callee)
		for _, arg := range t.arguments {
			err := r.expr_resolve(arg)
			if err != nil {
				return nil
			}
		}
	case *Get:
		err := r.expr_resolve(t.object)
		if err != nil {
			return err
		}
	case *Grouping:
		err := r.expr_resolve(t.expression)
		if err != nil {
			return nil
		}
	case *Literal:
		return nil
	case *Logical:
		err := r.expr_resolve(t.left)
		if err != nil {
			return nil
		}
		err = r.expr_resolve(t.right)
		if err != nil {
			return nil
		}
	case *Set:
		err := r.expr_resolve(t.value)
		if err != nil {
			return nil
		}
		err = r.expr_resolve(t.object)
		if err != nil {
			return nil
		}
	case *Unary:
		err := r.expr_resolve(t.right)
		if err != nil {
			return nil
		}
	case *Variable:
		if len(r.scopes) != 0 {
			front := r.scopes[len(r.scopes)-1]
			// Fun check to see if map[string]bool exists, and if it does, if the value is false
			// Extra fun around the default value of bools being false
			if v, ok := front[t.name.lexeme]; ok {
				if !v {
					tokenError(t.name, "Can't read local variable in its own initializer.")
				}
			}
		}
		r.resolveLocal(t, t.name)
	}
	return nil
}

func (r *Resolver) resolve_stmts(stmts []Stmt) error {
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

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1] // Pop a slice
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.lexeme]; ok {
		tokenError(name, "Already a variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true

}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if _, ok := scope[name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
