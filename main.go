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
		// Default to REPL if no command or file is passed
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: File path is required. Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: File path is required. Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Shorthand: check if first argument is a file with .rx extension
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
		} else {
			printHelp()
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start REPL")
	fmt.Println("  rustx repl             Start REPL")
	fmt.Println("  rustx run <file.rx>    Execute an .rx script")
	fmt.Println("  rustx build <file.rx>  Syntax check an .rx script")
	fmt.Println("  rustx version          Print compiler version")
	fmt.Println("  rustx <file.rx>        Shorthand to execute an .rx script")
}

func runFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error: Could not read file %s: %s\n", filepath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l, filepath)

	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	ev := evaluator.New(filepath)
	evaluated := ev.Eval(program, env)

	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintln(os.Stderr, errObj.Inspect())
			os.Exit(1)
		}
	}
}

func buildFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error: Could not read file %s: %s\n", filepath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l, filepath)

	p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: Syntax verified for %s\n", filepath)
}
