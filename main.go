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
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "version":
		fmt.Printf("RustX Starter version %s\n", VERSION)
	default:
		// If first arg is a file, run it (shorthand)
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version] [file.rx]")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}
