package main

import (
	"fmt"
	"os"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
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
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Println(VERSION)

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
		// If argument looks like a filename, try running it directly
		if len(command) > 3 && command[len(command)-3:] == ".rx" {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>")
			fmt.Println("  rustx build <file.rx>")
			fmt.Println("  rustx repl")
			fmt.Println("  rustx version")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if res != nil && res.Type() == object.ERROR_OBJ {
		fmt.Printf("RuntimeError: %s\n", res.Inspect())
		os.Exit(1)
	}

	// If a main() function was defined, execute it automatically
	if mainFnObj, ok := env.Get("main"); ok {
		if _, isFn := mainFnObj.(*object.Function); isFn {
			callMain := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainRes := evaluator.Eval(callMain, env)
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Printf("RuntimeError: %s\n", mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for %s:\n", filename)
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated %s\n", filename)
}
