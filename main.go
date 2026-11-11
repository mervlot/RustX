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
)

const Version = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Bare 'rustx' starts REPL
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Println(Version)

	case "help", "-h", "--help":
		printHelp()

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing file path for 'run' command.")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing file path for 'build' command.")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If command looks like a file (e.g. rustx script.rx), execute it directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command or file: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println(Version)
	fmt.Println("A lightweight hybrid scripting and systems language inspired by Rust and Go.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx                 Start the interactive REPL")
	fmt.Println("  rustx repl            Start the interactive REPL")
	fmt.Println("  rustx run <file.rx>   Execute a RustX source file")
	fmt.Println("  rustx <file.rx>       Execute a RustX source file directly")
	fmt.Println("  rustx build <file.rx> Parse and validate a RustX source file")
	fmt.Println("  rustx version         Display version information")
	fmt.Println("  rustx help            Display this help message")
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
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

	// If a main() function was declared, invoke it
	if mainFunc, ok := env.Get("main"); ok {
		if fn, isFunc := mainFunc.(*object.Function); isFunc && len(fn.Parameters) == 0 {
			callNode := &ast.CallExpression{
				Function: &ast.Identifier{Value: "main"},
			}
			mainRes := evaluator.Eval(callNode, env)
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Fprintln(os.Stderr, mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Validation passed: %s successfully parsed.\n", filename)
}
