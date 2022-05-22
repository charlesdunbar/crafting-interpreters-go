package main

type Stmt interface {
	Statement() Stmt
}

type Expression struct {
	expression Expr
}

type Var struct {
	name Token
	initializer Expr
}

type Print struct {
	expression Expr
}

func (e *Expression) Statement() Stmt { return e }

func (e *Var) Statement() Stmt { return e }

func (e *Print) Statement() Stmt { return e }

