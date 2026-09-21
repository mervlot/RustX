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

const VERSION = "RustX Starter v0.1"

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
		fmt.Println(VERSION)

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
		// If argument ends in .rx or file exists, treat as shorthand 'run'
		if len(args) >= 1 {
			if _, err := os.Stat(command); err == nil {
				runFile(command)
				return
			}
		}
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("RustX Starter")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script file")
	fmt.Println("  rustx <file.rx>       Shorthand to run a script file")
	fmt.Println("  rustx repl            Start interactive REPL prompt")
	fmt.Println("  rustx build <file.rx> Check script syntax and validate build")
	fmt.Println("  rustx version         Display version information")
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
		fmt.Printf("Syntax error(s) in %s:\n", filename)
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// Check if a main function is defined and execute it if so
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFunc := mainFn.(*object.Function); isFunc && len(fn.Parameters) == 0 {
			mainProgram := parser.New(lexer.New("main()")).ParseProgram()
			if len(mainProgram.Statements) > 0 {
				mainCallResult := evaluator.Eval(mainProgram.Statements[0], env)
				if mainCallResult != nil && mainCallResult.Type() == object.ERROR_OBJ {
					fmt.Println(mainCallResult.Inspect())
					os.Exit(1)
				}
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
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s parsed with 0 syntax errors.\n", filename)
}
