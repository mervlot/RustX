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
)

const VERSION = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Println(VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	case "help", "-h", "--help":
		printUsage()

	default:
		// If first argument is a file (e.g. `rustx hello.rx`), run it directly
		if fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println(VERSION)
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Execute a RustX script file")
	fmt.Println("  rustx repl           Start the interactive REPL shell")
	fmt.Println("  rustx build <file.rx> Check script for syntax and parsing errors")
	fmt.Println("  rustx version        Print RustX version")
	fmt.Println("  rustx <file.rx>      Directly run a script file")
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	input := string(data)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		fmt.Printf("Syntax error(s) in %s:\n", filename)
		for _, parseErr := range errs {
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

	// If a zero-argument `main` function is defined, invoke it automatically as entry point
	if mainFnObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainFnObj.(*object.Function); ok && len(mainFn.Parameters) == 0 {
			mainRes := evaluator.Eval(mainFn.Body, object.NewEnclosedEnvironment(mainFn.Env))
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Println(mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
	_ = ast.Node(nil)
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		fmt.Printf("Build failed for %s:\n", filename)
		for _, parseErr := range errs {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful! %s passed all lexical and syntax checks.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
