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
	fmt.Println("RustX Starter - A lightweight scripting + systems hybrid language")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  run <file.rx>     Execute a RustX source file")
	fmt.Println("  repl              Start the interactive REPL prompt")
	fmt.Println("  build <file.rx>   Validate and compile-check a RustX source file")
	fmt.Println("  version           Display version information")
	fmt.Println("  help              Show this help message")
	fmt.Println("\nShorthand:")
	fmt.Println("  rustx <file.rx>   Execute a RustX source file directly")
	fmt.Println("  rustx             Start REPL")
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
			fmt.Fprintln(os.Stderr, "Error: missing filename. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If argument ends with .rx or is an existing file, run it
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command or file: %s\n\n", command)
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
		fmt.Fprintf(os.Stderr, "Error: could not read file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
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

	// If main() function was declared, invoke main()
	if mainFunc, ok := env.Get("main"); ok {
		if fn, ok := mainFunc.(*object.Function); ok && len(fn.Parameters) == 0 {
			mainCall := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainResult := evaluator.Eval(mainCall, env)
			if mainResult != nil && mainResult.Type() == object.ERROR_OBJ {
				fmt.Fprintln(os.Stderr, mainResult.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for %s\n", filename)
}
