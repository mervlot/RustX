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
	braceDepth := 0
	parenDepth := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if len(buffer) == 0 && (trimmed == "exit" || trimmed == "quit") {
			break
		}

		buffer = append(buffer, line)

		// Calculate balance/depth
		for i := 0; i < len(line); i++ {
			ch := line[i]
			// Crude but robust tracking for a REPL
			// Skip character tracking inside string literals to avoid counting '{' inside strings
			if ch == '"' {
				// skip until next quote in the same line
				i++
				for i < len(line) && line[i] != '"' {
					i++
				}
				continue
			}

			// Single line comment
			if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				break
			}

			if ch == '{' {
				braceDepth++
			} else if ch == '}' {
				braceDepth--
			} else if ch == '(' {
				parenDepth++
			} else if ch == ')' {
				parenDepth--
			}
		}

		if braceDepth > 0 || parenDepth > 0 {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		// Balanced! Evaluate the accumulated input
		input := strings.Join(buffer, "\n")
		buffer = nil
		braceDepth = 0
		parenDepth = 0

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			if err, ok := evaluated.(*object.Error); ok {
				fmt.Fprintln(out, err.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
