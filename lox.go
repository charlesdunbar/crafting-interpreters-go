package main

import (
	"bufio"
	"fmt"
	"os"
)

var hadError bool
var hadRuntimeError bool

type Lox struct {

}

func (l *Lox) RunFile(source string) {
	f, err := os.ReadFile(source)
	if err != nil {
		panic("Error!")
	}
	l.run(string(f))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func (l *Lox) RunPrompt() {
	fmt.Print("> ")
	// Handles Ctrl-D for us
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l.run(s.Text())
		hadError = false
		fmt.Print("> ")
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens(l)
	parser := NewParser(tokens, l)
	statements := parser.parse()
	if hadError {
		return
	}
	err := NewInterpreter().interpret(statements)
	if err != nil {
		runtimeError(err)
	}
}

func lox_error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	hadError = true
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func tokenError(token Token, message string) {
	if token.l_type == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func runtimeError(e error) {
	fmt.Println(e)
	hadRuntimeError = true
}
