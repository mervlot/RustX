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

const VERSION = "0.1.0"

func printUsage() {
	fmt.Println("RustX Starter - Lightweight scripting & systems hybrid language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Execute a RustX source file")
	fmt.Println("  rustx build <file.rx> Validate and check a RustX source file")
	fmt.Println("  rustx repl            Start the interactive REPL prompt")
	fmt.Println("  rustx version         Show version information")
	fmt.Println("  rustx <file.rx>       Directly run a RustX source file")
	fmt.Println("  rustx                 Start the interactive REPL prompt (default)")
}

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

	case "help", "-h", "--help":
		printUsage()

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'run' command.")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'build' command.")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If argument looks like a .rx file or existing file, run it
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
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
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, result.Inspect())
		os.Exit(1)
	}

	// If a main() function was declared, execute it
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFn := mainFn.(*object.Function); isFn && len(fn.Parameters) == 0 {
			mainRes := evaluator.Eval(fn.Body, object.NewEnclosedEnvironment(fn.Env))
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Fprintln(os.Stderr, mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: '%s' syntax and structure verified.\n", filename)
}
