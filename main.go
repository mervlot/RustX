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
	"strings"
)

const VERSION = "0.1.0"

func printUsage() {
	fmt.Println("RustX Starter - Lightweight Scripting & Systems Hybrid Language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start interactive REPL")
	fmt.Println("  rustx repl             Start interactive REPL")
	fmt.Println("  rustx run <file.rx>    Execute a RustX source file")
	fmt.Println("  rustx <file.rx>        Execute a RustX source file (shorthand)")
	fmt.Println("  rustx build <file.rx>  Check and compile/validate a RustX source file")
	fmt.Println("  rustx version          Print RustX version")
	fmt.Println("  rustx help             Print this help message")
}

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

	case "help", "-h", "--help":
		printUsage()

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for build command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for run command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	default:
		// If argument ends with .rx or is a file, execute it directly
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.NewWithFilename(l, filename)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated '%s' (%d statements).\n", filename, len(program.Statements))
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.NewWithFilename(l, filename)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
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

	// If a main() function is defined in the script, call it automatically
	if mainFunc, ok := env.Get("main"); ok {
		if fn, ok := mainFunc.(*object.Function); ok && len(fn.Parameters) == 0 {
			res := evaluator.Eval(&ast.CallExpression{
				Function: &ast.Identifier{Value: "main"},
			}, env)
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Println(res.Inspect())
				os.Exit(1)
			}
		}
	}
}
