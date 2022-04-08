package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	cmdArgs := os.Args[1:]

	if len(cmdArgs) != 1 {
		fmt.Println("Usage: generateAst <output directory>")
		os.Exit(64)
	}
	outputDir := cmdArgs[0]
	println(outputDir)

	defineAst(outputDir, "Expr", []string{
		"Binary    : left Expr, operator Token, right Expr",
		"Grouping  : expression Expr",
		"Literal   : value interface{}",
		"Unary     : operator Token, right Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("package main\n\n")
	f.WriteString("type " + baseName + " interface {\n\tString() string\n}\n\n")
	f.WriteString("")
	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(f, baseName, className, fields)
	}
}

func defineType(writer *os.File, baseName, className, fieldList string) {
	writer.WriteString("type " + className + " struct {\n")
	//writer.WriteString("	Expr\n")
	fields := strings.Split(fieldList, ", ")
	for _, f := range fields {
		name := strings.Split(f, " ")[0]
		l_type := strings.Split(f, " ")[1]
		writer.WriteString("\t" + name + " " + l_type + "\n")
	}
	writer.WriteString("}\n\n")
}
