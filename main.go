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

const Version = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Shorthand for REPL when no arguments are provided
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter %s\n", Version)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	default:
		// Shorthand execution: if passing a valid .rx file path as the first argument,
		// execute it directly, bypassing the need for the 'run' command.
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n", cmd)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>")
			fmt.Println("  rustx <file.rx>      (Shorthand execution)")
			fmt.Println("  rustx repl")
			fmt.Println("  rustx version")
			os.Exit(1)
		}
	}
}

func runFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filepath, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		if errObj, ok := result.(*object.Error); ok {
			fmt.Println(errObj.Inspect())
			os.Exit(1)
		}
	}
}
