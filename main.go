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

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter %s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		path := args[1]
		if filepath.Ext(path) != ".rx" {
			fmt.Fprintln(os.Stderr, "Error: file must have .rx extension")
			os.Exit(1)
		}
		os.Exit(buildFile(path))

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		path := args[1]
		if filepath.Ext(path) != ".rx" {
			fmt.Fprintln(os.Stderr, "Error: file must have .rx extension")
			os.Exit(1)
		}
		os.Exit(runFile(path))

	default:
		if filepath.Ext(cmd) == ".rx" {
			os.Exit(runFile(cmd))
		}

		fmt.Fprintln(os.Stderr, "Unknown command. Available commands: run, build, repl, version")
		os.Exit(1)
	}
}

func runFile(path string) int {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", path, err.Error())
		return 1
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, errMsg := range p.Errors() {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		return 1
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, evaluated.Inspect())
		return 1
	}

	return 0
}

func buildFile(path string) int {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", path, err.Error())
		return 1
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, errMsg := range p.Errors() {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		return 1
	}

	fmt.Println("Build successful (syntax OK).")
	return 0
}
