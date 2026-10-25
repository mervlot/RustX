package main

import (
	"fmt"
	"os"
	"rustx/ast"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"rustx/repl"
	"rustx/token"
	"strings"
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
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version|build] <file.rx>")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	code := string(bytes)
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// If main function is defined in script, invoke main()
	if mainFnObj, ok := env.Get("main"); ok {
		if _, isFn := mainFnObj.(*object.Function); isFn {
			callMain := &ast.CallExpression{
				Token: token.Token{Type: token.LPAREN, Literal: "(", Line: 1, Column: 1},
				Function: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "main", Line: 1, Column: 1},
					Value: "main",
				},
				Arguments: []ast.Expression{},
			}
			mainCallResult := evaluator.Eval(callMain, env)
			if mainCallResult != nil && mainCallResult.Type() == object.ERROR_OBJ {
				fmt.Println(mainCallResult.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	code := string(bytes)
	l := lexer.New(code)
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		fmt.Printf("Build failed for %s\n", filename)
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for %s\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
