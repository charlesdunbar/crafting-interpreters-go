package main

import "fmt"

type Token struct {
	l_type  TokenType
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s, %s, %+v", t.l_type, t.lexeme, t.literal)
}
