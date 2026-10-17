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
		// Check if shorthand execution with .rx extension or valid file
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p, err := parser.New(l)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p, err := parser.New(l)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = p.ParseProgram()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	fmt.Printf("Build successful! Validated syntax for '%s'. Output program target: %s\n", filename, base)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printUsage() {
	fmt.Println("\nRustX Starter CLI")
	fmt.Println("Usage:")
	fmt.Println("  rustx                     Start REPL")
	fmt.Println("  rustx run <file.rx>       Run a script file")
	fmt.Println("  rustx build <file.rx>     Check and validate a script file")
	fmt.Println("  rustx version             Print RustX version")
	fmt.Println("  rustx <file.rx>           Shorthand run")
}
