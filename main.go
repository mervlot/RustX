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

const VERSION = "0.1.0"

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

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing file path for 'run'")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing file path for 'build'")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Shorthand: rustx <file.rx>
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			printUsage()
			os.Exit(1)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printUsage() {
	fmt.Println("RustX Starter - Lightweight Scripting & Systems Hybrid Language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start interactive REPL")
	fmt.Println("  rustx run <file.rx>    Execute a RustX program")
	fmt.Println("  rustx build <file.rx>  Validate and compile a RustX program")
	fmt.Println("  rustx repl             Start interactive REPL")
	fmt.Println("  rustx version          Print language version")
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if res != nil && res.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, res.Inspect())
		os.Exit(1)
	}

	// Check if a `main` function was declared and execute it if present
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			callNode := &ast.CallExpression{
				Token:    token.Token{Type: token.LPAREN, Literal: "("},
				Function: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "main"}, Value: "main"},
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
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s (%d statements parsed)\n", filename, len(program.Statements))
}
