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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Println("RustX Starter v0.1")
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: please provide a file path to run.")
			fmt.Fprintln(os.Stderr, "Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	default:
		// Check for shorthand file execution, e.g. rustx hello.rx
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read file %s: %s\n", path, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Fprintln(os.Stderr, parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintf(os.Stderr, "RuntimeError:\n%s\nat line %d column %d\n", errObj.Message, errObj.Token.Line, errObj.Token.Column)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start interactive REPL")
	fmt.Println("  rustx repl             Start interactive REPL")
	fmt.Println("  rustx run <file.rx>    Run an .rx script file")
	fmt.Println("  rustx <file.rx>        Shorthand to run an .rx script file")
	fmt.Println("  rustx version          Display version")
}
