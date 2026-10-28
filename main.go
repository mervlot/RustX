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
		// If argument ends with .rx or is a valid file, run directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file not found: %s\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Syntax Error(s):")
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if res != nil && res.Type() == object.ERROR_OBJ {
		fmt.Println(res.Inspect())
		os.Exit(1)
	}

	// Check if `main` function is defined and invoke it
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			mainCall := &ast.CallExpression{
				Token:     token.Token{Type: token.LPAREN, Literal: "("},
				Function:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "main"}, Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainRes := evaluator.Eval(mainCall, env)
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Println(mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed with %d syntax error(s):\n", len(p.Errors()))
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated and built '%s'\n", filename)
}

func printHelp() {
	fmt.Println("RustX Starter - Scripting and systems language hybrid")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println("  rustx <file.rx>")
	fmt.Println("\nCommands:")
	fmt.Println("  run <file.rx>   Execute a RustX script")
	fmt.Println("  repl            Start interactive REPL prompt")
	fmt.Println("  build <file.rx> Validate script syntax without execution")
	fmt.Println("  version         Show RustX version")
	fmt.Println("  help            Show this help message")
}
