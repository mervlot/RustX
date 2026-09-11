package main

import (
	"fmt"
	"os"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
)

const Version = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version", "-v", "--version":
		fmt.Println(Version)

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
		printUsage()

	default:
		// Shorthand execution: rustx file.rx
		if len(args) == 1 && len(command) > 0 && command[0] != '-' {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			printUsage()
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

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s parsed with no syntax errors.\n", filename)
}

func printUsage() {
	fmt.Println("RustX Starter - A lightweight scripting + systems hybrid language")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx run <file.rx>   Execute a RustX script")
	fmt.Println("  rustx <file.rx>       Execute a RustX script directly")
	fmt.Println("  rustx repl            Start the interactive REPL prompt")
	fmt.Println("  rustx build <file.rx> Validate syntax of a RustX script")
	fmt.Println("  rustx version         Show version information")
}
