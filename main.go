package main

import (
	"fmt"
	"os"

	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"rustx/repl"
)

const VERSION = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Println(VERSION)

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

	default:
		// If argument ends with .rx or file exists, treat as run <file>
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file not found: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>   Run a RustX script")
	fmt.Println("  rustx repl          Start interactive REPL")
	fmt.Println("  rustx build <file.rx> Check syntax of script without running")
	fmt.Println("  rustx version       Print version information")
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

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
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

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
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

	fmt.Printf("Build successful: '%s' syntax is valid.\n", filename)
}
