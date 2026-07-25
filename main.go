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

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

	case "run":
		if len(args) < 2 {
			fmt.Println("Error: 'run' command expects a file path.")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	default:
		if strings.HasSuffix(cmd, ".rx") {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n", cmd)
			fmt.Println("Usage:")
			fmt.Println("  rustx run <file.rx>      - Run a script")
			fmt.Println("  rustx <file.rx>          - Shorthand to run a script")
			fmt.Println("  rustx repl               - Start interactive REPL")
			fmt.Println("  rustx                    - Start interactive REPL")
			fmt.Println("  rustx version            - Display version info")
			os.Exit(1)
		}
	}
}

func runFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: failed to read file %s: %s\n", filePath, err.Error())
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
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintln(os.Stderr, "RuntimeError:")
			fmt.Fprintln(os.Stderr, errObj.Message)
			os.Exit(1)
		}
	}
}
