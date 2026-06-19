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

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := os.Args[1]

	switch command {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <filename>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	default:
		// Assume run if a file is provided without 'run' command
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("Parser error: %s\n", msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Printf("Runtime error: %s\n", evaluated.Inspect())
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx              Starts REPL")
	fmt.Println("  rustx repl         Starts REPL")
	fmt.Println("  rustx run <file>   Runs a .rx file")
	fmt.Println("  rustx version      Shows version")
}
