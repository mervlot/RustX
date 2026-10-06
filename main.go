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

const VERSION = "RustX Starter v0.1"

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
		fmt.Println(VERSION)
	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// If first arg ends with .rx or is a file, run it directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Commands: run <file.rx>, repl, version, build <file.rx>")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
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

	if errObj, ok := result.(*object.Error); ok {
		fmt.Println(errObj.Inspect())
		os.Exit(1)
	}

	// Check if a `main` function is defined in the script and call it automatically
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn {
			if len(fn.Parameters) == 0 {
				evaluator.Eval(fn.Body, env)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	_ = p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for '%s':\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build check successful for '%s'. No syntax errors found.\n", filename)
}
