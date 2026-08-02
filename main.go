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

const VERSION = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Printf("RustX Starter %s\n", VERSION)
		return

	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return

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
		// Shorthand execution: check if it's a file path ending in .rx
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
			return
		}

		fmt.Printf("Unknown command or file: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
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

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}

func buildFile(path string) {
	ext := filepath.Ext(path)
	if ext != ".rx" {
		fmt.Printf("Error: file %s is not a .rx file\n", path)
		os.Exit(1)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Println("Build successful (syntax verified).")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  rustx                  Start REPL")
	fmt.Println("  rustx repl             Start REPL")
	fmt.Println("  rustx run <file.rx>    Run the specified .rx script")
	fmt.Println("  rustx <file.rx>        Shorthand to run the specified .rx script")
	fmt.Println("  rustx build <file.rx>  Lexically and syntactically validate the file")
	fmt.Println("  rustx version          Print version info")
}
