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

const VERSION = "0.1.0"

func printUsage() {
	fmt.Println("RustX Starter - Lightweight Scripting & Systems Hybrid Language")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx run <file.rx>     Execute a RustX source file")
	fmt.Println("  rustx build <file.rx>   Check syntax and validate file")
	fmt.Println("  rustx repl              Start interactive REPL session")
	fmt.Println("  rustx version           Display version information")
	fmt.Println("  rustx <file.rx>         Shorthand to execute a file")
	fmt.Println("  rustx                   Start REPL session")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "help", "-h", "--help":
		printUsage()

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Shorthand: rustx file.rx
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, result.Inspect())
		os.Exit(1)
	}

	// Check if a zero-argument `main` function is defined and execute it
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFunc := mainFn.(*object.Function); isFunc && len(fn.Parameters) == 0 {
			callNode := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			callResult := evaluator.Eval(callNode, env)
			if callResult != nil && callResult.Type() == object.ERROR_OBJ {
				fmt.Fprintln(os.Stderr, callResult.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s verified without errors.\n", filename)
}
