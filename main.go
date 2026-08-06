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

const VERSION = "0.1"

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

	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		path := args[1]
		errs := validateFile(path)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			os.Exit(1)
		}
		fmt.Printf("Build successful: %s contains no syntax errors.\n", path)
		return

	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		path := args[1]
		executeFile(path)
		return

	default:
		if filepath.Ext(cmd) == ".rx" || fileExists(cmd) {
			executeFile(cmd)
			return
		}

		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Usage:")
		fmt.Println("  rustx                  Launch REPL")
		fmt.Println("  rustx repl             Launch REPL")
		fmt.Println("  rustx version          Show version")
		fmt.Println("  rustx build <file.rx>  Validate syntax")
		fmt.Println("  rustx run <file.rx>    Execute file")
		fmt.Println("  rustx <file.rx>        Execute file (shorthand)")
		os.Exit(1)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func validateFile(path string) []error {
	content, err := os.ReadFile(path)
	if err != nil {
		return []error{fmt.Errorf("error reading file %s: %w", path, err)}
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	return p.Errors()
}

func executeFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
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
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, result.Inspect())
		os.Exit(1)
	}
}
