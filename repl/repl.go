package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/username/rustx/errors"
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
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")

	var inputBuffer strings.Builder
	braceCount := 0
	parenCount := 0

	fmt.Fprint(out, PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" || line == "quit" {
			break
		}

		inputBuffer.WriteString(line + "\n")

		// Count braces and parens ignoring strings
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == '"' {
				// Skip string contents
				i++
				for i < len(line) && line[i] != '"' {
					if line[i] == '\\' {
						i++
					}
					i++
				}
				continue
			}
			if ch == '{' {
				braceCount++
			} else if ch == '}' {
				braceCount--
			} else if ch == '(' {
				parenCount++
			} else if ch == ')' {
				parenCount--
			}
		}

		if braceCount > 0 || parenCount > 0 {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		// Reset counters for next evaluation
		braceCount = 0
		parenCount = 0

		fullInput := inputBuffer.String()
		inputBuffer.Reset()

		if strings.TrimSpace(fullInput) == "" {
			fmt.Fprint(out, PROMPT)
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if err, ok := evaluated.(*object.Error); ok {
				fmt.Fprintln(out, err.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
