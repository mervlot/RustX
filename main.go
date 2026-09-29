package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/repl"
)

const VERSION = "0.1.0"

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
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)
	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: rustx run <filename.rx>")
			os.Exit(1)
		}
		runFile(args[1])
	case "build":
		if len(args) < 2 {
			fmt.Println("Usage: rustx build <filename.rx>")
			os.Exit(1)
		}
		buildFile(args[1])
	default:
		// If argument ends with .rx or file exists, execute directly: rustx file.rx
		if strings.HasSuffix(command, ".rx") || fileExists(command) {
			runFile(command)
		} else {
			fmt.Printf("Unknown command or file: %s\n", command)
			fmt.Println("Usage: rustx [run|repl|version|build] <file.rx>")
			os.Exit(1)
		}
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printErrors(p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// Auto-invoke main() if defined in script and zero parameters
	if mainObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainObj.(*object.Function); ok && len(mainFn.Parameters) == 0 {
			callExpr := &ast.CallExpression{
				Function:  &ast.Identifier{Value: "main"},
				Arguments: []ast.Expression{},
			}
			res := evaluator.Eval(callExpr, env)
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Println(res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}

	input := string(bytes)
	l := lexer.New(input)
	p := parser.New(l)

	_ = p.ParseProgram()
	if len(p.Errors()) != 0 {
		printErrors(p.Errors())
		os.Exit(1)
	}

	fmt.Printf("Build successful: %s contains no syntax errors.\n", filename)
}

func printErrors(errors []error) {
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
