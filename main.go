package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	case "version":
		fmt.Println(VERSION)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		if fileExists(command) {
			runFile(command)
		} else {
			printUsage()
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
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// If a main() function was defined, execute main()
	if mainFn, ok := env.Get("main"); ok {
		if fn, isFn := mainFn.(*object.Function); isFn {
			if len(fn.Parameters) == 0 {
				mainRes := evaluator.Eval(fn.Body, object.NewEnclosedEnvironment(fn.Env))
				if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
					fmt.Println(mainRes.Inspect())
					os.Exit(1)
				}
			}
		}
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

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: '%s' has no syntax errors.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() && filepath.Ext(filename) == ".rx"
}

func printUsage() {
	fmt.Println(VERSION)
	fmt.Println("Usage:")
	fmt.Println("  rustx run <file.rx>    Run a script file")
	fmt.Println("  rustx repl             Start interactive REPL prompt")
	fmt.Println("  rustx build <file.rx>  Validate syntax without running")
	fmt.Println("  rustx version          Print language version")
}
