package main

import (
	"fmt"
	"os"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
)

const Version = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Println(Version)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
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
		// Check if argument is a file path
		if len(args) == 1 && fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("RustX Starter Compiler / Interpreter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script")
	fmt.Println("  rustx <file.rx>       Run a RustX script directly")
	fmt.Println("  rustx repl            Start interactive REPL prompt")
	fmt.Println("  rustx build <file.rx> Parse and validate syntax of a file")
	fmt.Println("  rustx version         Show version information")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

	if len(p.Errors()) != 0 {
		fmt.Printf("Syntax Error(s) in %s:\n", filename)
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

	// Auto-invoke main() in script file if defined and no error occurred
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn {
			mainRes := evaluator.ApplyFunction(nil, fn, []object.Object{})
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

	if len(p.Errors()) != 0 {
		fmt.Printf("Build failed with syntax error(s) in %s:\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
		}
		os.Exit(1)
	}

	fmt.Printf("Build succeeded for %s (0 syntax errors).\n", filename)
}
