package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
)

const Version = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Printf("RustX Starter %s\n", Version)

	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <filepath.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <filepath.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If arg ends with .rx or file exists, run script directly as shorthand
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|build|repl|version] [file.rx]")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
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

	// Auto-invoke main function if registered in environment
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn {
			if len(fn.Parameters) == 0 {
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
}

func buildFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
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

	fmt.Printf("Successfully validated syntax for '%s'. Ready for compilation/execution.\n", filename)
}
