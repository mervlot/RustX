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
		startRepl()
		return
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "repl":
		startRepl()
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	default:
		// Default to run if it looks like a file, or show help
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Println("Unknown command:", command)
			fmt.Println("Usage: rustx [run|repl|version] [file.rx]")
		}
	}
}

func startRepl() {
	fmt.Printf("RustX Starter v%s\n", VERSION)
	repl.Start(os.Stdin, os.Stdout)
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Println("SyntaxError:")
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
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
