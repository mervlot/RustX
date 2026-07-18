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
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Println("RustX Starter v0.1")

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(os.Args[2])

	default:
		// Shorthand direct file execution: rustx <file.rx>
		if strings.HasSuffix(command, ".rx") {
			runFile(command)
			return
		}

		// Fallback print usage
		fmt.Println("RustX Starter CLI")
		fmt.Println("Commands:")
		fmt.Println("  rustx run <file.rx>   Execute a RustX script")
		fmt.Println("  rustx repl            Start interactive REPL prompt")
		fmt.Println("  rustx version         Show version info")
		fmt.Println("\nShorthand:")
		fmt.Println("  rustx <file.rx>       Directly run a .rx script file")
		os.Exit(1)
	}
}

func runFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file %q: %s\n", filepath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Compilation failed due to syntax errors:")
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
