package main

import (
	"bytes"
	"fmt"
	"go/format"
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
		"Assign    : name Token, value Expr",
		"Binary    : left Expr, operator Token, right Expr",
		"Call	   : callee Expr, paren Token, arguments []Expr",
		"Get       : object Expr, name Token",
		"Grouping  : expression Expr",
		"Literal   : value any",
		"Logical   : left Expr, operator Token, right Expr",
		"Set       : object Expr, name Token, value Expr",
		"Super     : keyword Token, method Token",
		"This      : keyword Token",
		"Unary     : operator Token, right Expr",
		"Variable  : name Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block      : statements []Stmt",
		"Class		: name Token, superclass Variable, methods []Function",
		"Expression : expression Expr",
		"Function   : name Token, params []Token, body []Stmt",
		"If         : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Var        : name Token, initializer Expr",
		"Print      : expression Expr",
		"Return     : keyword Token, value Expr",
		"While      : condition Expr, body Stmt",
	})
}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := new(bytes.Buffer)

	buf.WriteString("package main\n\n")
	var fun string
	if baseName == "Expr" {
		fun = "Expression()"
	} else if baseName == "Stmt" {
		fun = "Statement()"
	}
	buf.WriteString("type " + baseName + " interface {\n\t" + fun + " " + baseName + "\n}\n\n")
	buf.WriteString("")
	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(buf, baseName, className, fields)
	}

	// Fake method to make the types not be any
	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		returnSelf(buf, baseName, className, fun)
	}

	content, _ := format.Source(buf.Bytes())
	f.Write(content)
}

func defineType(writer *bytes.Buffer, baseName, className, fieldList string) {
	writer.WriteString("type " + className + " struct {\n")
	fields := strings.Split(fieldList, ", ")
	for _, f := range fields {
		name := strings.Split(f, " ")[0]
		l_type := strings.Split(f, " ")[1]
		writer.WriteString("\t" + name + " " + l_type + "\n")
	}
	writer.WriteString("}\n\n")
}

func returnSelf(writer *bytes.Buffer, baseName, className, fun string) {
	writer.WriteString("func (e *" + className + ") " + fun + " " + baseName + " { return e }\n\n")
}
