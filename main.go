package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
	"github.com/username/rustx/token"
)

const Version = "RustX Starter v0.1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "version", "-v", "--version":
		fmt.Println(Version)
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
		// Direct script execution (e.g. rustx script.rx)
		if filepath.Ext(cmd) == ".rx" || fileExists(cmd) {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command or file: %s\n", cmd)
			fmt.Println("Usage: rustx [run|repl|version|build] <file.rx>")
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

	// Auto-invoke main() if defined in top-level environment
	if mainFn, ok := env.Get("main"); ok {
		if _, isFn := mainFn.(*object.Function); isFn {
			mainCallNode := &ast.CallExpression{
				Token:    token.Token{Type: token.LPAREN, Literal: "(", Line: 1, Column: 1},
				Function: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "main", Line: 1, Column: 1}, Value: "main"},
				Arguments: []ast.Expression{},
			}
			mainRes := evaluator.Eval(mainCallNode, env)
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

	if len(p.Errors()) > 0 {
		fmt.Printf("Build failed with %d syntax error(s):\n", len(p.Errors()))
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s checked with zero syntax errors.\n", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
