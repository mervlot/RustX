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
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "help":
		printHelp()
	default:
		// Check if the first arg is a file ending in .rx
		if len(os.Args) == 2 && isRxFile(os.Args[1]) {
			runFile(os.Args[1])
		} else {
			printHelp()
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
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}

func isRxFile(path string) bool {
	return len(path) > 3 && path[len(path)-3:] == ".rx"
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>  - Execute a .rx script")
	fmt.Println("  rustx repl           - Start interactive REPL")
	fmt.Println("  rustx version        - Show version")
	fmt.Println("  rustx <file.rx>      - Shorthand for run")
}
