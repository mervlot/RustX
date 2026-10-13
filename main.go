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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If argument ends in .rx or is a file that exists, run it directly
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version|build] <file.rx>")
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

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}

	// Check if a zero-argument `main` function is defined and execute it
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			callExpr := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainResult := evaluator.Eval(callExpr, env)
			if mainResult != nil && mainResult.Type() == object.ERROR_OBJ {
				fmt.Println(mainResult.Inspect())
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
	p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully built %s (Syntax clean)\n", filePath)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
