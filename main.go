package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
)

const VERSION = "v0.1.0"

func main() {
	args := os.Args[1:]

	// If no arguments, or the command is repl, run REPL
	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]
	switch cmd {
	case "version":
		fmt.Printf("RustX Starter %s\n", VERSION)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: run command requires a file path argument. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	default:
		// Shorthand support: if passing a valid .rx file directly, execute it
		if filepath.Ext(cmd) == ".rx" {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command: %s. Use 'repl', 'run <file.rx>', or 'version'.\n", cmd)
			os.Exit(1)
		}
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", path, err.Error())
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
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, evaluated.Inspect())
		os.Exit(1)
	}
}
