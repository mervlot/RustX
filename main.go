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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Shorthand for REPL
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]
	switch cmd {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: No file specified for 'run' command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: No file specified for 'build' command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		// Shorthand direct file execution: if argument is a file with .rx extension
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n", cmd)
			fmt.Println("Usage:")
			fmt.Println("  rustx                (Launch REPL)")
			fmt.Println("  rustx repl           (Launch REPL)")
			fmt.Println("  rustx run <file.rx>  (Run a script)")
			fmt.Println("  rustx build <file.rx> (Build/verify a script)")
			fmt.Println("  rustx <file.rx>      (Shorthand direct run)")
			fmt.Println("  rustx version        (Print version)")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: Could not read file %s\n%v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
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
		fmt.Printf("Error: Could not read file %s\n%v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful (Syntax Check Passed) for %s\n", filename)
}
