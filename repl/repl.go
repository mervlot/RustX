package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
)

const PROMPT = ">>> "
const MULTI_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type commands or Ctrl+D to exit\n")

	var buffer []string

	for {
		if len(buffer) == 0 {
			fmt.Fprintf(out, PROMPT)
		} else {
			fmt.Fprintf(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		buffer = append(buffer, line)

		fullInput := strings.Join(buffer, "\n")
		if isBalanced(fullInput) {
			l := lexer.New(fullInput)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				buffer = nil
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintf(out, "%s\n", errObj.Inspect())
				} else {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				}
			}
			buffer = nil
		}
	}
}

func isBalanced(input string) bool {
	// Simple tracker for balance of curly braces and parentheses
	// Ignore contents inside strings if possible, but for simplicity:
	braceCount := 0
	parenCount := 0
	inString := false
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				inString = false
			}
		} else {
			if ch == '"' {
				inString = true
			} else if ch == '{' {
				braceCount++
			} else if ch == '}' {
				braceCount--
			} else if ch == '(' {
				parenCount++
			} else if ch == ')' {
				parenCount--
			}
		}
	}

	return braceCount <= 0 && parenCount <= 0 && !inString
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
