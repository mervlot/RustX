package main

import (
	"fmt"
	"os"
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
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]
	switch command {
	case "version":
		fmt.Println("RustX Starter v0.1")
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: 'run' command requires a file path.")
			os.Exit(1)
		}
		runFile(args[1])
		return

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: 'build' command requires a file path.")
			os.Exit(1)
		}
		buildFile(args[1])
		return

	default:
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
			return
		}

		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Launch the interactive REPL")
	fmt.Println("  rustx run <file.rx>    Execute a RustX Starter script")
	fmt.Println("  rustx build <file.rx>  Syntactically validate a script without execution")
	fmt.Println("  rustx version          Print version info")
	fmt.Println("  rustx <file.rx>        Execute a script directly (shorthand)")
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %q: %v\n", path, err)
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
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Println(errObj.Inspect())
			os.Exit(1)
		}
	}
}

func buildFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %q: %v\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Println("Validation OK")
}
