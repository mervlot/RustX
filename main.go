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

const VERSION = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	cmd := args[0]

	switch cmd {
	case "version", "-v", "--version":
		fmt.Printf("RustX Starter v%s\n", VERSION)

	case "repl":
		repl.Start(os.Stdin, os.Stdout)

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

	case "help", "-h", "--help":
		printHelp()

	default:
		// If argument ends with .rx or is a file, run directly
		if strings.HasSuffix(cmd, ".rx") || fileExists(cmd) {
			runFile(cmd)
		} else {
			fmt.Printf("Unknown command: %s\n\n", cmd)
			printHelp()
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("RustX Starter Compiler / Interpreter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rustx <command> [arguments]")
	fmt.Println("  rustx <file.rx>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file.rx>   Execute a RustX script file")
	fmt.Println("  repl            Start the interactive REPL shell")
	fmt.Println("  build <file.rx> Validate and check syntax of a script file")
	fmt.Println("  version         Display version info")
	fmt.Println("  help            Display this help message")
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}

	// If main() function exists in top-level env, invoke it
	if mainFunc, ok := env.Get("main"); ok {
		if fn, isFunc := mainFunc.(*object.Function); isFunc && len(fn.Parameters) == 0 {
			mainCall := &astCallHelper{fn: fn}
			res := mainCall.invoke(env)
			if res != nil && res.Type() == object.ERROR_OBJ {
				fmt.Println(res.Inspect())
				os.Exit(1)
			}
		}
	}
}

func buildFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("Build failed with %d syntax error(s):\n", len(p.Errors()))
		for _, parseErr := range p.Errors() {
			fmt.Println(parseErr.Error())
		}
		os.Exit(1)
	}

	fmt.Printf("Successfully validated %s - no syntax errors found.\n", filename)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type astCallHelper struct {
	fn *object.Function
}

func (h *astCallHelper) invoke(env *object.Environment) object.Object {
	bodyEnv := object.NewEnclosedEnvironment(h.fn.Env)
	res := evaluator.Eval(h.fn.Body, bodyEnv)
	if ret, ok := res.(*object.ReturnValue); ok {
		return ret.Value
	}
	return res
}
