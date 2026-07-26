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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Println("RustX Starter v0.1")
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
		return

	default:
		if filepath.Ext(cmd) == ".rx" {
			runFile(cmd)
			return
		}

		fmt.Printf("Unknown command or file: %s\n", cmd)
		fmt.Println("Usage:")
		fmt.Println("  rustx run <file.rx>")
		fmt.Println("  rustx repl")
		fmt.Println("  rustx version")
		fmt.Println("  rustx <file.rx>  (shorthand)")
		os.Exit(1)
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %q: %s\n", path, err)
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
