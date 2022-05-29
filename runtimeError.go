package main

import "fmt"

type RuntimeError struct {
	token Token
	err   string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		token: token,
		err:   message,
	}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf(e.err + "\n[line " + fmt.Sprint(e.token.line) + "]")
}
