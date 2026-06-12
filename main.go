package main

import (
	"fmt"
	"io"
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
		startREPL()
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
		startREPL()
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	default:
		// Try to run the file if no command is specified but an argument is given
		runFile(os.Args[1])
	}
}

func startREPL() {
	fmt.Printf("RustX Starter v%s\n", VERSION)
	fmt.Println("Type 'exit' or press Ctrl+C to quit.")
	repl.Start(os.Stdin, os.Stdout)
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if errObj, ok := evaluated.(*object.Error); ok {
		fmt.Fprintf(os.Stderr, "RuntimeError: %s\n", errObj.Message)
		os.Exit(1)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintf(out, "SyntaxError:\n")
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
