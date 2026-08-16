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

const Version = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter %s\n", Version)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

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

	case "help", "-h", "--help":
		printHelp()

	default:
		// If arg is a file path ending in .rx or existing file, run directly
		if _, err := os.Stat(cmd); err == nil {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", cmd)
			printHelp()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
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

	// Check if main() function exists and call it
	if mainObj, ok := env.Get("main"); ok {
		if _, ok := mainObj.(*object.Function); ok {
			callMain := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			callMainRes := evaluator.Eval(callMain, env)
			if callMainRes != nil && callMainRes.Type() == object.ERROR_OBJ {
				fmt.Println(callMainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for %s\n", filename)
}

func printHelp() {
	fmt.Println("RustX Starter Language CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx [command] [file]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.rx>   Execute a RustX script")
	fmt.Println("  build <file.rx> Validate and check syntax of a script")
	fmt.Println("  repl            Start interactive REPL prompt")
	fmt.Println("  version         Display language version")
}
