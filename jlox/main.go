package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lox := Lox{false}
	cmdArgs := os.Args[1:]

	if len(cmdArgs) == 0 {
		lox.runPrompt()
	} else if len(cmdArgs) == 1 {
		lox.runFile(cmdArgs[0])
	} else {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64)
	}
}

type Lox struct {
	hadError bool
}

func (l *Lox) runFile(source string) {
	f, err := os.ReadFile(source)
	if err != nil {
		panic("Error!")
	}
	l.run(string(f))
	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) runPrompt() {
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
	tokens := scanner.scanTokens()
	for _, t := range tokens {
		fmt.Println(t)
	}
}
