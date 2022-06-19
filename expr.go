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

type Variable struct {
	name Token
}

type Unary struct {
	operator Token
	right    Expr
}

func (e *Assign) Expression() Expr { return e }

func (e *Binary) Expression() Expr { return e }

func (e *Grouping) Expression() Expr { return e }

func (e *Literal) Expression() Expr { return e }

func (e *Logical) Expression() Expr { return e }

func (e *Variable) Expression() Expr { return e }

func (e *Unary) Expression() Expr { return e }
