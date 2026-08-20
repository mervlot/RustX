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
		// If arg ends with .rx or file exists, run file directly
		if filepath.Ext(command) == ".rx" || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.EvalProgram(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, evaluated.Inspect())
		os.Exit(1)
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed for %s:\n", filename)
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build succeeded for %s: Syntax OK.\n", filename)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx                 Start the interactive REPL prompt")
	fmt.Println("  rustx run <file.rx>   Execute a RustX source file")
	fmt.Println("  rustx <file.rx>       Shorthand to execute a RustX source file")
	fmt.Println("  rustx build <file.rx> Validate syntax of a RustX source file")
	fmt.Println("  rustx version         Display version information")
}
