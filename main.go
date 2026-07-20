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

const VERSION = "0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	switch args[0] {
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: File path is required for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
		return

	default:
		// Shorthand execution: check if it's a file ending in .rx
		filePath := args[0]
		if strings.HasSuffix(filePath, ".rx") {
			runFile(filePath)
			return
		}

		fmt.Printf("Unknown command or file: %s\n", filePath)
		fmt.Println("Usage:")
		fmt.Println("  rustx                  (Starts REPL)")
		fmt.Println("  rustx repl             (Starts REPL)")
		fmt.Println("  rustx run <file.rx>    (Runs a script)")
		fmt.Println("  rustx <file.rx>        (Shorthand to run a script)")
		fmt.Println("  rustx version          (Prints version)")
		os.Exit(1)
	}
}

func runFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: Could not read file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Syntax Errors found:")
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
			os.Exit(1)
		}
	}
}
