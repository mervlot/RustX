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
	"strings"
)

const Version = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", Version)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	case "help", "-h", "--help":
		printHelp()
	default:
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx                     Start interactive REPL")
	fmt.Println("  rustx run <file.rx>       Execute a RustX source file")
	fmt.Println("  rustx <file.rx>           Direct shortcut to run a script file")
	fmt.Println("  rustx build <file.rx>     Check and validate a RustX source file for syntax errors")
	fmt.Println("  rustx repl                Start interactive REPL")
	fmt.Println("  rustx version             Print RustX version")
	fmt.Println("  rustx help                Display this help message")
}

func runFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Errors in '%s':\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if errObj, ok := res.(*object.Error); ok {
		fmt.Println(errObj.Inspect())
		os.Exit(1)
	}

	// Check if a `main` function is defined and execute it automatically
	if mainFnObj, ok := env.Get("main"); ok {
		if fn, ok := mainFnObj.(*object.Function); ok && len(fn.Parameters) == 0 {
			callExpr := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainRes := evaluator.Eval(callExpr, env)
			if errObj, ok := mainRes.(*object.Error); ok {
				fmt.Println(errObj.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for '%s':\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: '%s' passed syntax and lexical check.\n", filename)
}
