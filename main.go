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

const Version = "v0.1.0"

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
		fmt.Printf("RustX Starter %s\n", Version)

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
		// Shorthand: rustx <file.rx>
		if len(args) == 1 && (fileExists(args[0]) || len(args[0]) > 3 && args[0][len(args[0])-3:] == ".rx") {
			runFile(args[0])
		} else {
			printUsage()
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(code))
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

	// Auto-invoke main() if defined in the script's global scope
	if mainObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainObj.(*object.Function); ok {
			mainRes := evaluator.ApplyFunction(mainFn, []object.Object{}, token.Token{})
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Println(mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(code))
	p := parser.New(l)

	_ = p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful for %s! Syntax and lexical analysis passed.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printUsage() {
	fmt.Println("RustX Starter - Lightweight Systems + Scripting Language")
	fmt.Println("\nUsage:")
	fmt.Println("  rustx run <file.rx>   Run a .rx script file")
	fmt.Println("  rustx <file.rx>       Shorthand execution of a .rx script")
	fmt.Println("  rustx repl            Start interactive REPL prompt")
	fmt.Println("  rustx build <file.rx> Validate syntax of a .rx file")
	fmt.Println("  rustx version         Show version information")
}
