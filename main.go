package main

import (
	"fmt"
	"os"
	"strings"

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
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: 'run' command requires a filepath. E.g., 'rustx run hello.rx'")
			os.Exit(1)
		}
		runFile(args[1])
		return

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: 'build' command requires a filepath. E.g., 'rustx build hello.rx'")
			os.Exit(1)
		}
		buildFile(args[1])
		return

	default:
		// Shorthand execution: if passing a valid .rx file directly
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
			return
		}

		fmt.Printf("Unknown command or invalid file: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx               - Start the interactive REPL")
	fmt.Println("  rustx repl          - Start the interactive REPL")
	fmt.Println("  rustx run <file.rx> - Run the specified script")
	fmt.Println("  rustx <file.rx>     - Run the specified script (shorthand)")
	fmt.Println("  rustx build <file>  - Syntactic and lexical check of the file")
	fmt.Println("  rustx version       - Print version information")
}

func runFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filepath, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Fprintln(os.Stderr, evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func buildFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filepath, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Println("Build successful (syntax check passed).")
}
