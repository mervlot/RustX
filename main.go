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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// No arguments: launch REPL by default
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: 'run' command requires a file path argument")
			os.Exit(1)
		}
		runFile(args[1])
		return

	default:
		// Check for shorthand execution: passing a valid .rx file directly
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
			return
		}

		// Invalid usage
		fmt.Printf("Unknown command or file: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filepath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintln(os.Stderr, "SyntaxError:")
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "  %s\n", msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintln(os.Stderr, errObj.Inspect())
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx               - Start the interactive REPL")
	fmt.Println("  rustx repl          - Start the interactive REPL")
	fmt.Println("  rustx run <file>    - Run the specified .rx script")
	fmt.Println("  rustx <file.rx>     - Shorthand to run the specified .rx script directly")
	fmt.Println("  rustx version       - Print version information")
}
