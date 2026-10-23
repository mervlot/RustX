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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

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
			fmt.Println("Error: missing script file path. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If argument ends with .rx or is a valid file path, default to run file
		if _, err := os.Stat(cmd); err == nil {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file '%s'\n", cmd)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(program, env)

	if res != nil && res.Type() == object.ERROR_OBJ {
		fmt.Println(res.Inspect())
		os.Exit(1)
	}

	// If a main() function is defined, invoke it as the entry point
	if mainObj, ok := env.Get("main"); ok {
		if fnObj, isFn := mainObj.(*object.Function); isFn && len(fnObj.Parameters) == 0 {
			mainCall := &ast.CallExpression{
				Token:    token.Token{Type: token.LPAREN, Literal: "("},
				Function: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "main"}, Value: "main"},
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
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("Build failed for %s:\n", filename)
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s parsed with 0 errors.\n", filename)
}

func printUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script")
	fmt.Println("  rustx <file.rx>       Run a RustX script directly")
	fmt.Println("  rustx repl            Start interactive REPL")
	fmt.Println("  rustx build <file.rx> Check syntax of script")
	fmt.Println("  rustx version         Show version info")
}
