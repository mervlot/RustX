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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

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
		// Shorthand direct run: check if argument ends with .rx or is a valid file
		if filepath.Ext(cmd) == ".rx" {
			runFile(cmd)
		} else {
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("RustX CLI usage:")
	fmt.Println("  rustx                         Start interactive REPL")
	fmt.Println("  rustx repl                    Start interactive REPL")
	fmt.Println("  rustx run <file.rx>           Run specified .rx file")
	fmt.Println("  rustx <file.rx>               Shorthand execution of .rx file")
	fmt.Println("  rustx build <file.rx>         Perform syntactic validation of .rx file")
	fmt.Println("  rustx version                 Display the language version")
}

func runFile(filename string) {
	if filepath.Ext(filename) != ".rx" {
		fmt.Println("Error: File must have .rx extension")
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Syntax errors found:")
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintln(os.Stderr, errObj.Message)
			os.Exit(1)
		}
	}
}

func buildFile(filename string) {
	if filepath.Ext(filename) != ".rx" {
		fmt.Println("Error: File must have .rx extension")
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Build failed with syntax errors:")
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	fmt.Println("Build succeeded. Syntax is valid.")
}
