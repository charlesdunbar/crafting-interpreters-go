package main

import (
	"fmt"
	"os"

	. "github.com/charlesdunbar/lox-go/scanner"
)

func main() {
	lox := Lox{}
	cmdArgs := os.Args[1:]

	if len(cmdArgs) == 0 {
		lox.RunPrompt()
	} else if len(cmdArgs) == 1 {
		lox.RunFile(cmdArgs[0])
	} else {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	}
}
