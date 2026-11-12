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

func printUsage() {
	fmt.Println("RustX Starter - A lightweight scripting + systems hybrid language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script")
	fmt.Println("  rustx build <file.rx> Validate/check syntax of a RustX script")
	fmt.Println("  rustx repl            Start interactive REPL prompt")
	fmt.Println("  rustx version         Show version information")
	fmt.Println("  rustx <file.rx>       Shorthand to run a script directly")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Default to REPL if no arguments provided
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", Version)
	case "help", "-h", "--help":
		printUsage()
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: missing filename for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing filename for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	default:
		// Shorthand: rustx file.rx
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n\n", command)
			printUsage()
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

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
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

	// Check if a zero-argument main() function exists and invoke it
	if mainFnObj, ok := env.Get("main"); ok {
		if fn, isFn := mainFnObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			callResult := evaluator.Eval(&ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}, env)
			if callResult != nil && callResult.Type() == object.ERROR_OBJ {
				fmt.Println(callResult.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Compiled and validated '%s' successfully.\n", filename)
}
