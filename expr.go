package main

type Expr interface {
	String() string
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value interface{}
}

type Unary struct {
	operator Token
	right Expr
}

