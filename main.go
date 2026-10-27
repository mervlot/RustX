package main

import (
	"fmt"
	"os"
	"path/filepath"

	"rustx/ast"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"rustx/repl"
)

const VERSION = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Println(VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Check if argument is a file path
		if filepath.Ext(command) == ".rx" || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>")
			fmt.Println("  rustx build <file.rx>")
			fmt.Println("  rustx repl")
			fmt.Println("  rustx version")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if errObj, ok := res.(*object.Error); ok {
		fmt.Printf("RuntimeError: %s\n", errObj.Message)
		os.Exit(1)
	}

	// Check if main() function exists and call it
	if mainObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainObj.(*object.Function); ok {
			callExpr := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			_ = mainFn
			mainRes := evaluator.Eval(callExpr, env)
			if errObj, ok := mainRes.(*object.Error); ok {
				fmt.Printf("RuntimeError in main(): %s\n", errObj.Message)
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s (Syntax and lexical check passed)\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
