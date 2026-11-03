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

const VERSION = "RustX Starter v0.1"

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Execute a RustX program")
	fmt.Println("  rustx <file.rx>       Execute a RustX program (shorthand)")
	fmt.Println("  rustx build <file.rx> Check syntax of a RustX program")
	fmt.Println("  rustx repl            Start interactive REPL")
	fmt.Println("  rustx                 Start interactive REPL (default)")
	fmt.Println("  rustx version         Show version information")
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
		fmt.Println(VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "help", "-h", "--help":
		printUsage()

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: missing file argument for build command")
			printUsage()
			os.Exit(1)
		}
		buildFile(args[1])

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing file argument for run command")
			printUsage()
			os.Exit(1)
		}
		runFile(args[1])

	default:
		// Check if argument ends with .rx or is a file path
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
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

	// Check if a `main` function was declared and execute it if present
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFunc := mainFn.(*object.Function); isFunc {
			if len(fn.Parameters) == 0 {
				callExpr := &ast.CallExpression{
					Token:    fn.Body.Token,
					Function: &ast.Identifier{Token: fn.Body.Token, Value: "main"},
				}
				mainResult := evaluator.Eval(callExpr, env)
				if mainResult != nil && mainResult.Type() == object.ERROR_OBJ {
					fmt.Fprintln(os.Stderr, mainResult.Inspect())
					os.Exit(1)
				}
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
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

	fmt.Printf("Build check succeeded for '%s': 0 errors found.\n", filename)
}
