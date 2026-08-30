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
			fmt.Println("Error: Missing file argument for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file argument for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// Shorthand: rustx file.rx
		if len(args) == 1 && (len(command) > 3 && command[len(command)-3:] == ".rx") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|build|repl|version] [file.rx]")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// If main() function is defined, auto-execute it
	if mainFnObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainFnObj.(*object.Function); ok {
			evaluator.Eval(&ast.CallExpression{
				Function:  &ast.Identifier{Token: mainFn.Body.Token, Value: "main"},
				Arguments: []ast.Expression{},
			}, env)
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Build check successful for %s (no syntax errors)\n", filename)
}
