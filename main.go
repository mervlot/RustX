package main

import (
	"fmt"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
	"os"
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
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: No file specified for 'run' command.")
			os.Exit(1)
		}
		runFile(os.Args[2])
	default:
		// If it's not a command, maybe it's a file path to run directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			showHelp()
			os.Exit(1)
		}
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
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
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  rustx repl          - Start the interactive REPL")
	fmt.Println("  rustx run <file>    - Run a .rx script")
	fmt.Println("  rustx <file>        - Run a .rx script directly")
	fmt.Println("  rustx version       - Show version information")
}
