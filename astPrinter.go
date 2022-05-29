package main

import "fmt"

func (b Binary) String() string {
	return parenthesize(b.operator.lexeme, b.left, b.right)
}

func (g Grouping) String() string {
	return parenthesize("group", g.expression)
}

func (l Literal) String() string {
	if l.value == nil {
		return "nil"
	}
	return fmt.Sprint(l.value)
}

func (u Unary) String() string {
	return parenthesize(u.operator.lexeme, u.right)
}

func parenthesize(name string, exprs ...Expr) string {
	s := ""
	s += "(" + name
	for _, e := range exprs {
		s += " "
		s += fmt.Sprint(e)
	}
	s += ")"

	return s
}

func test() {
	u := Unary{}
	u.operator = Token{MINUS, "-", nil, 1}
	u.right = &Literal{value: 123}

	// x := Unary{
	// 	operator: Token{MINUS, "-", nil, 1},
	// 	right:    Literal{value: 123},
	// }
	expression := Binary{
		&Unary{
			Token{MINUS, "-", nil, 1},
			&Literal{123},
		},
		Token{STAR, "*", nil, 1},
		&Grouping{expression: &Literal{45.67}},
	}
	fmt.Println(expression)
}
