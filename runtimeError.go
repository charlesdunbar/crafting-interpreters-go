package main

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
	return e.err
}
