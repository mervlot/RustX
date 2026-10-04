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

const Version = "v0.1.0"

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

	case "version":
		fmt.Printf("RustX Starter %s\n", Version)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// If arg ends with .rx or is a file, execute directly
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start interactive REPL")
	fmt.Println("  rustx run <file.rx>    Execute a RustX script file")
	fmt.Println("  rustx build <file.rx>  Validate syntax of a RustX file")
	fmt.Println("  rustx version          Display version")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Printf("%s\n", err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Printf("%s\n", result.Inspect())
		os.Exit(1)
	}

	// Check if a zero-argument main function exists in env and invoke it
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFn := mainFn.(*object.Function); isFn && len(fn.Parameters) == 0 {
			res := evaluator.Eval(fn.Body, object.NewEnclosedEnvironment(fn.Env))
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Printf("%s\n", res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for %s with %d syntax error(s):\n", filename, len(p.Errors()))
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err.Error())
		}
		os.Exit(1)
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("Successfully validated syntax for %s\n", absPath)
}
