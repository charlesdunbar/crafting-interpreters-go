package main

type Stmt interface {
	Statement() Stmt
}

type Block struct {
	statements []Stmt
}

type Expression struct {
	expression Expr
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type Var struct {
	name        Token
	initializer Expr
}

type Print struct {
	expression Expr
}

func (e *Block) Statement() Stmt { return e }

func (e *Expression) Statement() Stmt { return e }

func (e *If) Statement() Stmt { return e }

func (e *Var) Statement() Stmt { return e }

func (e *Print) Statement() Stmt { return e }
