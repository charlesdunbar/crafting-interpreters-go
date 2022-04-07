package main

type Expr struct {}

type Binary struct {
	Expr
	left Expr
	operator Token
	right Expr
}
type Grouping struct {
	Expr
	expression Expr
}
type Literal struct {
	Expr
	value interface{}
}
type Unary struct {
	Expr
	operator Token
	right Expr
}
