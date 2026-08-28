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

	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <filepath.rx>")
			os.Exit(1)
		}
		runFile(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <filepath.rx>")
			os.Exit(1)
		}
		buildFile(args[1])

	default:
		if fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file not found: %s\n", command)
			fmt.Println("Commands:")
			fmt.Println("  rustx run <file.rx>   Run a RustX script")
			fmt.Println("  rustx build <file.rx> Check and validate a RustX script")
			fmt.Println("  rustx repl            Start interactive REPL")
			fmt.Println("  rustx version         Show version")
			os.Exit(1)
		}
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
		fmt.Printf("Compilation failed with %d error(s):\n", len(p.Errors()))
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Printf("RuntimeError: %s\n", evaluated.Inspect())
		os.Exit(1)
	}

	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			res := evaluator.ApplyFunction(fn, []object.Object{})
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Printf("RuntimeError: %s\n", res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file %q: %s\n", filepath, err.Error())
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed with %d error(s):\n", len(p.Errors()))
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully built and validated %s\n", filepath)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
