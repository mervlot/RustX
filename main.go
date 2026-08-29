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
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", Version)

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
		// If first argument is a file path ending with .rx or valid file, run directly
		if fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>")
			fmt.Println("  rustx repl")
			fmt.Println("  rustx build <file.rx>")
			fmt.Println("  rustx version")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	codeBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	code := string(codeBytes)
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// Auto-execute zero-arg main() function if defined at top-level upon script execution
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			callExp := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainResult := evaluator.Eval(callExp, env)
			if mainResult != nil && mainResult.Type() == object.ERROR_OBJ {
				fmt.Println(mainResult.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	codeBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	code := string(codeBytes)
	l := lexer.New(code)
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Build failed with syntax errors:")
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
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
