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

const VERSION = "v0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// Default to repl
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version":
		fmt.Printf("RustX Starter %s\n", VERSION)
		return
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
		return
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: Please specify a file path to run.\nUsage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
		return
	default:
		// Check if shorthand execute for a valid .rx file path
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
			return
		}
		// Invalid command
		fmt.Printf("Error: Unknown command or invalid file '%s'\n", cmd)
		fmt.Println("Usage:\n  rustx run <file.rx>   Execute a script\n  rustx repl            Start interactive REPL\n  rustx version         Show version info\n  rustx <file.rx>       Shorthand execution")
		os.Exit(1)
	}
}

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("SyntaxError:")
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
			os.Exit(1)
		}
	}
}
