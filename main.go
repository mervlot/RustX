package main

import (
	"fmt"
	"os"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"rustx/repl"
	"strings"
)

const Version = "0.1.0"

func printUsage() {
	fmt.Println("RustX Starter - Lightweight scripting & systems programming language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.rx>    Execute a RustX source file")
	fmt.Println("  repl             Start the interactive REPL prompt")
	fmt.Println("  version          Print the version information")
	fmt.Println("  build <file.rx>  Check and compile/validate a RustX file")
	fmt.Println("  <file.rx>        Shorthand to run a RustX source file")
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

	case "version", "--version", "-v":
		fmt.Printf("RustX Starter v%s\n", Version)

	case "help", "--help", "-h":
		printUsage()

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'run' command")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'build' command")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Check if argument is a .rx file shorthand
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Syntax errors in %s:\n", filename)
		for _, parseErr := range p.Errors() {
			parseErr.File = filename
			fmt.Fprintf(os.Stderr, "%s\n", parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, result.Inspect())
		os.Exit(1)
	}

	// If a main function is defined, invoke it
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFunc := mainFn.(*object.Function); isFunc {
			res := evaluator.Eval(fn.Body, object.NewEnclosedEnvironment(fn.Env))
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Fprintln(os.Stderr, res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed with %d error(s):\n", len(p.Errors()))
		for _, parseErr := range p.Errors() {
			parseErr.File = filename
			fmt.Fprintf(os.Stderr, "%s\n", parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build check successful for %s: no syntax errors detected.\n", filename)
}
