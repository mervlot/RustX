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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Press Ctrl+D or Ctrl+C to exit.")

	var buffer []string

	for {
		if len(buffer) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		buffer = append(buffer, line)

		joined := strings.Join(buffer, "\n")
		if isBalanced(joined) {
			l := lexer.New(joined)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
			buffer = nil // Reset buffer
		}
	}
}

// isBalanced tracks balanced braces and parentheses for multiline inputs
func isBalanced(input string) bool {
	var braces int
	var parens int
	var inString bool
	var escaped bool

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		switch ch {
		case '{':
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		}
	}

	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
