package main

import (
	"fmt"
	"os"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
	"github.com/username/rustx/token"
)

const VERSION = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := args[0]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter %s\n", VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for build command.")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		filename := args[1]
		err := validateFileSyntax(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Build successful: %s parsed without errors.\n", filename)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Missing file path for run command.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		filename := args[1]
		executeFile(filename)

	default:
		// If first argument is a file name, execute it directly
		if fileExists(command) {
			executeFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version|build] [file.rx]")
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

func validateFileSyntax(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", filename, err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("%s", p.Errors()[0])
	}

	return nil
}

func executeFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintf(os.Stderr, "%s\n", result.Inspect())
		os.Exit(1)
	}

	// Automatic entry-point invocation if main() function exists at top-level
	if mainFunc, ok := env.Get("main"); ok {
		if fn, ok := mainFunc.(*object.Function); ok && len(fn.Parameters) == 0 {
			mainRes := evaluator.Eval(&astCallMain{fn: fn}, env)
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Fprintf(os.Stderr, "%s\n", mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

type astCallMain struct {
	fn *object.Function
}

func (a *astCallMain) TokenLiteral() string { return "main" }
func (a *astCallMain) String() string       { return "main()" }
func (a *astCallMain) GetToken() token.Token { return token.Token{} }
