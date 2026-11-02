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
	fmt.Println("RustX Starter - A lightweight scripting + systems hybrid language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start the interactive REPL")
	fmt.Println("  rustx repl             Start the interactive REPL")
	fmt.Println("  rustx run <file.rx>    Execute a RustX source file")
	fmt.Println("  rustx build <file.rx>  Check and compile/validate a RustX source file")
	fmt.Println("  rustx version          Display version information")
	fmt.Println("  rustx help             Show this help message")
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

	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "help", "--help", "-h":
		printUsage()

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'run'")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: missing filename for 'build'")
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If argument ends in .rx or is a file, run it directly
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
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

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	res := evaluator.Eval(prog, env)

	if res != nil && res.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, res.Inspect())
		os.Exit(1)
	}

	// If main() function was defined, call it as entry point if top-level did not already return an error
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
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	// Successful validation and syntax check
	stmtCount := len(prog.Statements)
	_ = ast.Program{}
	fmt.Printf("Build successful: %s parsed correctly (%d top-level statements).\n", filename, stmtCount)
}
