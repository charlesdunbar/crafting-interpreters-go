package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoxTestScripts(t *testing.T) {
	files, err := os.ReadDir("./lox_programs/")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == "fail.lox" ||
			file.Name() == "resolver-errors.lox" {
			continue
		}
		a := NewLox()
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Printf("==== Running test %s ====\n", file.Name())
		a.RunFile(filepath.Join(path, "lox_programs/", file.Name()))
	}
}
