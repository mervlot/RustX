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
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>     Execute a .rx script file")
	fmt.Println("  rustx <file.rx>         Execute a .rx script file directly")
	fmt.Println("  rustx repl              Start the interactive REPL prompt")
	fmt.Println("  rustx                   Start the interactive REPL prompt")
	fmt.Println("  rustx build <file.rx>   Syntactically validate a .rx file")
	fmt.Println("  rustx version           Display version information")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "--version", "-v":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "help", "--help", "-h":
		printUsage()
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'run' command")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'build' command")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// If argument ends with .rx or file exists, treat as run <file>
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
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
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

	// If main() function was declared, automatically invoke it
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFn := mainFn.(*object.Function); isFn && len(fn.Parameters) == 0 {
			mainCall := &ast.CallExpression{
				Function: &ast.Identifier{Value: "main"},
			}
			mainRes := evaluator.Eval(mainCall, env)
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
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Validation successful: %s built without errors.\n", filename)
}
