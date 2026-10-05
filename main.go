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
	case "version", "--version", "-v":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
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
	case "help", "--help", "-h":
		printHelp()
	default:
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file '%s'\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("RustX Starter - Lightweight scripting + systems hybrid language")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx                     Start interactive REPL")
	fmt.Println("  rustx run <file.rx>       Execute a RustX source file")
	fmt.Println("  rustx build <file.rx>     Check syntax and validate a RustX source file")
	fmt.Println("  rustx repl                Start interactive REPL")
	fmt.Println("  rustx version             Display version information")
	fmt.Println("  rustx help                Display help menu")
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: Could not read file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}

	// Check if a zero-argument main function exists and call it
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			callExpr := &ast.CallExpression{
				Function:  &ast.Identifier{Tok: fn.Name.Tok, Value: fn.Name.Value},
				Arguments: []ast.Expression{},
			}
			res := evaluator.Eval(callExpr, env)
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Println(res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: Could not read file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("Build failed with %d syntax errors:\n", len(p.Errors()))
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: '%s' validated with 0 syntax errors.\n", filename)
}
