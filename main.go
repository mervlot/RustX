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
		// No arguments, start REPL
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
			fmt.Println("Error: run command expects a file path. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	default:
		// Check if the argument is a file with .rx extension
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>  (Execute a file)")
			fmt.Println("  rustx <file.rx>      (Execute a file directly)")
			fmt.Println("  rustx repl           (Start REPL)")
			fmt.Println("  rustx                (Start REPL)")
			fmt.Println("  rustx version        (Show version)")
			os.Exit(1)
		}
	}
}

func runFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error: could not read file %s: %s\n", filepath, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
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
