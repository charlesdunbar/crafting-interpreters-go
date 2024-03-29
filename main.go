package main

import (
	"fmt"
	"os"
)

func main() {
	lox := NewLox()
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
