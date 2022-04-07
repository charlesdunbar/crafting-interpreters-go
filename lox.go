package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
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
}

func (l *Lox) RunPrompt() {
	fmt.Print("> ")
	// Handles Ctrl-D for us
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Println(s.Text())
		l.run(s.Text())
		l.hadError = false
		fmt.Print("> ")
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens(l)
	for _, t := range tokens {
		fmt.Println(t)
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
