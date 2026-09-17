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
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Printf("RustX Starter %s\n", VERSION)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file argument for 'run' command.")
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		filePath := args[1]
		runFile(filePath)

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file argument for 'build' command.")
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		filePath := args[1]
		buildFile(filePath)

	case "help", "-h", "--help":
		printHelp()

	default:
		// If first argument is a filename ending in .rx or valid file path
		if fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func runFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filePath, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
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

	// Check if main() function was declared, and invoke it if zero arguments
	if mainFunc, ok := env.Get("main"); ok {
		if fn, ok := mainFunc.(*object.Function); ok && len(fn.Parameters) == 0 {
			mainIdent := &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "main", Line: 1, Column: 1},
				Value: "main",
			}
			callExpr := &ast.CallExpression{
				Token:     token.Token{Type: token.LPAREN, Literal: "(", Line: 1, Column: 1},
				Function:  mainIdent,
				Arguments: []ast.Expression{},
			}
			callRes := evaluator.Eval(callExpr, env)
			if callRes != nil && callRes.Type() == object.ERROR_OBJ {
				fmt.Println(callRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filePath, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for %s:\n", filePath)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for %s\n", filePath)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printHelp() {
	fmt.Println("RustX Starter Language CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script file")
	fmt.Println("  rustx repl          Start interactive REPL")
	fmt.Println("  rustx build <file.rx> Validate/check syntax of script")
	fmt.Println("  rustx version       Display version info")
	fmt.Println("  rustx help          Show this help message")
}
