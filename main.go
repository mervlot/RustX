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

const VERSION = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Default to REPL if no command or file is passed
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version":
		fmt.Println(VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: 'run' command requires a file path argument.")
			os.Exit(1)
		}
		filePath := args[1]
		runFile(filePath)
		return

	default:
		// If the argument ends with .rx or appears to be a file, run it directly
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
			return
		}

		// Otherwise, show usage help
		printUsage()
		os.Exit(1)
	}
}

func runFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %s\n", filePath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content), filePath)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "Syntax Error(s) in %s:\n", filePath)
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Fprintf(os.Stderr, "%s\n", evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx                     - Start REPL")
	fmt.Println("  rustx repl                - Start REPL")
	fmt.Println("  rustx run <file.rx>       - Execute a RustX script file")
	fmt.Println("  rustx <file.rx>           - Shorthand to execute a RustX script file")
	fmt.Println("  rustx version             - Print version information")
}
