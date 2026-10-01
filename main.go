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
	"github.com/username/rustx/token"
)

const VERSION = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter %s\n", VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing script file path. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: missing file path to build. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	case "help", "-h", "--help":
		printHelp()

	default:
		// If argument ends in .rx or exists as a file, shorthand run
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
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

	// If zero-argument main() is defined in environment, execute it as entrypoint
	if mainObj, ok := env.Get("main"); ok {
		if _, isFn := mainObj.(*object.Function); isFn {
			mainCall := &ast.ExpressionStatement{
				Expression: &ast.CallExpression{
					Token:    token.Token{Type: token.LPAREN, Literal: "("},
					Function: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "main"}, Value: "main"},
					Arguments: []ast.Expression{},
				},
			}
			res := evaluator.Eval(mainCall, env)
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Println(res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("Build failed with %d error(s):\n", len(p.Errors()))
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s parsed without errors.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printHelp() {
	fmt.Println("RustX Starter - Lightweight scripting & systems hybrid language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println("  rustx <file.rx>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.rx>   Execute a RustX source file")
	fmt.Println("  repl            Start interactive REPL session")
	fmt.Println("  build <file.rx> Validate and check syntax of a RustX source file")
	fmt.Println("  version         Display version info")
	fmt.Println("  help            Show this help message")
}
