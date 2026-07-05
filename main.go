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

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		startREPL()
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		printHelp()
	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "repl":
		startREPL()
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	default:
		// If first arg is a file, run it directly
		if _, err := os.Stat(command); err == nil {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printHelp()
			os.Exit(1)
		}
	}
}

func startREPL() {
	fmt.Printf("RustX Starter v%s\n", VERSION)
	repl.Start(os.Stdin, os.Stdout)
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if errObj, ok := evaluated.(*object.Error); ok {
		fmt.Println(errObj.Inspect())
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("RustX Starter CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx [command] [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  run <file.rx>  Execute a .rx script")
	fmt.Println("  repl           Start the interactive REPL")
	fmt.Println("  version        Show version information")
	fmt.Println("  help           Show this help message")
}
