package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
	"github.com/username/rustx/token"
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
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

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
		if strings.HasSuffix(command, ".rx") || filepath.Ext(command) == ".rx" {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("\nCommands:")
	fmt.Println("  rustx run <file.rx>   Execute a RustX script")
	fmt.Println("  rustx <file.rx>       Shorthand script execution")
	fmt.Println("  rustx repl            Start interactive REPL prompt")
	fmt.Println("  rustx build <file.rx> Perform syntax check and validate source")
	fmt.Println("  rustx version         Show version information")
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Failed to parse file '%s':\n", filename)
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

	// If main function is defined in script, invoke it automatically as entry point
	if mainFuncObj, ok := env.Get("main"); ok {
		if _, isFunc := mainFuncObj.(*object.Function); isFunc {
			mainCall := &ast.CallExpression{
				Tok: token.Token{Type: token.LPAREN, Literal: "("},
				Function: &ast.Identifier{
					Tok:   token.Token{Type: token.IDENT, Literal: "main"},
					Value: "main",
				},
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
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for '%s':\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for '%s'. Build OK.\n", filename)
}
