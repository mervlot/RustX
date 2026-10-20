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

const Version = "0.1.0"

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
		fmt.Printf("RustX Starter v%s\n", Version)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// If argument ends with .rx or file exists, run it directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version|build] [file.rx]")
			os.Exit(1)
		}
	}
}

func runFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filepath, err)
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

	if errObj, ok := evaluated.(*object.Error); ok {
		fmt.Println(errObj.Inspect())
		os.Exit(1)
	}

	// Check if a zero-argument main() function exists in the environment and run it
	if mainObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainObj.(*object.Function); ok {
			if len(mainFn.Parameters) == 0 {
				mainRes := evaluator.Eval(mainFn.Body, object.NewEnclosedEnvironment(mainFn.Env))
				if mainErr, ok := mainRes.(*object.Error); ok {
					fmt.Println(mainErr.Inspect())
					os.Exit(1)
				}
			}
		}
	}
}

func buildFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filepath, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for %s:\n", filepath)
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s parsed with 0 errors.\n", filepath)
}
