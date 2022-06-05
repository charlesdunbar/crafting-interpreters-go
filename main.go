package main

import (
	"fmt"
	"os"
)

func main() {
	lox := Lox{false, false, interpreter{
		environment: &Environment{
			values:    make(map[string]any),
			enclosing: nil,
		}}}
	cmdArgs := os.Args[1:]

	if len(cmdArgs) == 0 {
		//test()
		lox.RunPrompt()
	} else if len(cmdArgs) == 1 {
		lox.RunFile(cmdArgs[0])
	} else {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	}
}
