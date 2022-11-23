package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
}

func (l *Lox) RunFile(source string) {
	f, err := os.ReadFile(source)
	if err != nil {
		panic("Error!")
	}
	l.run(string(f))
	if l.hadError {
		os.Exit(65)
	}
	if l.hadRuntimeError {
		os.Exit(70)
	}
}

func (l *Lox) RunPrompt() {
	fmt.Print("> ")
	// Handles Ctrl-D for us
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l.run(s.Text())
		l.hadError = false
		fmt.Print("> ")
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens(l)
	parser := NewParser(tokens, l)
	statements := parser.parse()
	if l.hadError {
		return
	}
	err := NewInterpreter().interpret(statements)
	if err != nil {
		l.runtimeError(err)
	}
}

func (l Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	l.hadError = true
	//return fmt.Errorf("[line %d] Error%s: %s", line, where, message)
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func (l *Lox) tokenError(token Token, message string) {
	if token.l_type == EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (l *Lox) runtimeError(e error) {
	fmt.Println(e)
	l.hadRuntimeError = true
}
