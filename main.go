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
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "help":
		printHelp()
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	default:
		// Shorthand: rustx file.rx
		if isFile(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>  - Execute a RustX file")
	fmt.Println("  rustx repl           - Start interactive REPL")
	fmt.Println("  rustx version        - Show version info")
	fmt.Println("  rustx help           - Show this help")
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
	evaluator.Eval(program, env)
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
