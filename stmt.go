package main

type Stmt interface {
	Statement() Stmt
}

type Block struct {
	statements []Stmt
}

type Break struct {
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

type While struct {
	condition Expr
	body      Stmt
}

func (e *Block) Statement() Stmt { return e }

func (e *Break) Statement() Stmt { return e }

func (e *Expression) Statement() Stmt { return e }

func (e *If) Statement() Stmt { return e }

func (e *Var) Statement() Stmt { return e }

func (e *Print) Statement() Stmt { return e }

func (e *While) Statement() Stmt { return e }
