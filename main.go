package main

import (
	"fmt"
	"os"
	"github.com/username/rustx/repl"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/object"
)

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		printHelp()
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	default:
		// Shorthand: rustx file.rx
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printHelp()
			os.Exit(1)
		}
	}
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
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if errObj, ok := evaluated.(*object.Error); ok {
		fmt.Println(errObj.Inspect())
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx               - Start REPL")
	fmt.Println("  rustx repl          - Start REPL")
	fmt.Println("  rustx run <file>    - Run .rx script")
	fmt.Println("  rustx <file>        - Run .rx script (shorthand)")
	fmt.Println("  rustx version       - Show version")
	fmt.Println("  rustx help          - Show this help")
}
