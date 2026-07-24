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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()
		buffer = append(buffer, line)

		combined := strings.Join(buffer, "\n")
		if isBalanced(combined) {
			l := lexer.New(combined)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated != evaluator.NULL {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
			buffer = nil
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}
	}
}

func isBalanced(input string) bool {
	var braces int
	var parens int
	var inString bool
	var escape bool

	i := 0
	for i < len(input) {
		ch := input[i]

		if escape {
			escape = false
			i++
			continue
		}

		if ch == '\\' {
			escape = true
			i++
			continue
		}

		if ch == '"' {
			inString = !inString
			i++
			continue
		}

		if !inString {
			// Skip single-line comments: //
			if ch == '/' && i+1 < len(input) && input[i+1] == '/' {
				// skip until end of line
				for i < len(input) && input[i] != '\n' {
					i++
				}
				continue
			}

			// Skip multi-line comments: /* ... */
			if ch == '/' && i+1 < len(input) && input[i+1] == '*' {
				i += 2 // skip /*
				for i < len(input) {
					if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
						i += 2 // skip */
						break
					}
					i++
				}
				continue
			}

			if ch == '{' {
				braces++
			} else if ch == '}' {
				braces--
			} else if ch == '(' {
				parens++
			} else if ch == ')' {
				parens--
			}
		}

		i++
	}

	return braces <= 0 && parens <= 0 && !inString
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
