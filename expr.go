package main

type Expr interface {
	Expression() Expr
}

type Assign struct {
	name  Token
	value Expr
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

type Get struct {
	object Expr
	name   Token
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

type Set struct {
	object Expr
	name   Token
	value  Expr
}

type Super struct {
	keyword Token
	method  Token
}

type This struct {
	keyword Token
}

type Unary struct {
	operator Token
	right    Expr
}

type Variable struct {
	name Token
}

func (e *Assign) Expression() Expr { return e }

func (e *Binary) Expression() Expr { return e }

func (e *Call) Expression() Expr { return e }

func (e *Get) Expression() Expr { return e }

func (e *Grouping) Expression() Expr { return e }

func (e *Literal) Expression() Expr { return e }

func (e *Logical) Expression() Expr { return e }

func (e *Set) Expression() Expr { return e }

func (e *Super) Expression() Expr { return e }

func (e *This) Expression() Expr { return e }

func (e *Unary) Expression() Expr { return e }

func (e *Variable) Expression() Expr { return e }
