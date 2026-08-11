package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		// Run REPL by default
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version":
		fmt.Println("RustX Starter v0.1")

	case "build":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: 'build' command requires a file path (e.g., 'rustx build main.rx')")
			os.Exit(1)
		}
		filePath := args[1]
		if !strings.HasSuffix(filePath, ".rx") {
			fmt.Fprintln(os.Stderr, "Error: File must have .rx extension")
			os.Exit(1)
		}

		errorsFound := runBuild(filePath)
		if errorsFound {
			os.Exit(1)
		}
		fmt.Println("Build successful")
		os.Exit(0)

	case "run":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: 'run' command requires a file path (e.g., 'rustx run hello.rx')")
			os.Exit(1)
		}
		filePath := args[1]
		if !strings.HasSuffix(filePath, ".rx") {
			fmt.Fprintln(os.Stderr, "Error: File must have .rx extension")
			os.Exit(1)
		}
		runFile(filePath)

	default:
		// Check if the argument is a valid path to a .rx file for shorthand execution
		if strings.HasSuffix(command, ".rx") {
			_, err := os.Stat(command)
			if err == nil {
				runFile(command)
				return
			}
		}

		// Otherwise print usage instructions
		fmt.Fprintln(os.Stderr, "RustX Starter CLI")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  rustx                        Start the interactive REPL")
		fmt.Fprintln(os.Stderr, "  rustx repl                   Start the interactive REPL")
		fmt.Fprintln(os.Stderr, "  rustx version                Show version")
		fmt.Fprintln(os.Stderr, "  rustx run <file.rx>          Execute a .rx file")
		fmt.Fprintln(os.Stderr, "  rustx build <file.rx>        Perform syntax validation on a .rx file")
		fmt.Fprintln(os.Stderr, "  rustx <file.rx>              Shorthand file execution")
		os.Exit(1)
	}
}

func runFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not read file %s: %s\n", filePath, err)
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
	result := evaluator.Eval(program, env)

	if errObj, ok := result.(*object.Error); ok {
		fmt.Fprintln(os.Stderr, errObj.Inspect())
		os.Exit(1)
	}
}

func runBuild(filePath string) bool {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid path: %s\n", err)
		return true
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not read file %s: %s\n", filePath, err)
		return true
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return true
	}

	return false
}
