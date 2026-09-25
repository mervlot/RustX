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

const Version = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", Version)
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing script filename for 'run'")
			fmt.Println("Usage: rustx run <file.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Error: missing script filename for 'build'")
			fmt.Println("Usage: rustx build <file.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// If argument looks like a .rx file or exists, run it directly
		if strings.HasSuffix(cmd, ".rx") || fileExists(cmd) {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n", cmd)
			fmt.Println("Usage: rustx [run|repl|version|build] <file.rx>")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Printf("Syntax errors in %s:\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Printf("  %s\n", parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Printf("%s\n", evaluated.Inspect())
		os.Exit(1)
	}

	// If main() function is defined, invoke it
	if mainObj, ok := env.Get("main"); ok {
		if fn, ok := mainObj.(*object.Function); ok {
			mainRes := evaluator.Eval(fn.Body, object.NewEnclosedEnvironment(fn.Env))
			if mainRes != nil && mainRes.Type() == object.ERROR_OBJ {
				fmt.Printf("%s\n", mainRes.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	_ = p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed with syntax errors in %s:\n", filename)
		for _, parseErr := range p.Errors() {
			fmt.Printf("  %s\n", parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s verified with 0 syntax errors.\n", filepath.Base(filename))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
