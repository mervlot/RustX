package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If the argument is a valid file ending with .rx or existing file, run it directly
		if filepath.Ext(command) == ".rx" || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			printUsage()
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

	// Check if main() function exists and invoke it
	if mainFnObj, ok := env.Get("main"); ok {
		if mainFn, isFn := mainFnObj.(*object.Function); isFn && len(mainFn.Parameters) == 0 {
			callExpr := &ast.CallExpression{
				Function: &ast.Identifier{Value: "main"},
			}
			evalResult := evaluator.Eval(callExpr, env)
			if evalResult != nil && evalResult.Type() == object.ERROR_OBJ {
				fmt.Println(evalResult.Inspect())
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
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
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

func printUsage() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script")
	fmt.Println("  rustx repl          Start interactive REPL")
	fmt.Println("  rustx build <file.rx> Syntactically check/build script")
	fmt.Println("  rustx version       Display current version")
	fmt.Println("  rustx <file.rx>     Shorthand script execution")
}
