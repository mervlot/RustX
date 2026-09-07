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
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	case "help", "-h", "--help":
		printHelp()
	default:
		// If argument ends with .rx or file exists, execute directly
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, evaluated.Inspect())
		os.Exit(1)
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated syntax for '%s'. Ready for target compilation.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printHelp() {
	fmt.Println("RustX Starter - Lightweight hybrid scripting + systems programming language")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println("  rustx <file.rx>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.rx>   Execute a RustX source script file")
	fmt.Println("  repl            Start interactive RustX REPL session")
	fmt.Println("  build <file.rx> Validate and check syntax of a RustX source script")
	fmt.Println("  version         Display version info")
	fmt.Println("  help            Show this help message")
}
