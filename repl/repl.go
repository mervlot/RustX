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
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if openBraces == 0 && openParens == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if openBraces == 0 && openParens == 0 && trimmed == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Count open braces/parens ignoring comments or quotes simply
		inQuote := false
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c == '"' && (i == 0 || line[i-1] != '\\') {
				inQuote = !inQuote
			}
			if !inQuote {
				if c == '{' {
					openBraces++
				} else if c == '}' && openBraces > 0 {
					openBraces--
				} else if c == '(' {
					openParens++
				} else if c == ')' && openParens > 0 {
					openParens--
				}
			}
		}

		if openBraces == 0 && openParens == 0 {
			code := buffer.String()
			buffer.Reset()

			l := lexer.New(code)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				for _, err := range p.Errors() {
					fmt.Fprintf(out, "%s\n", err.Error())
				}
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				} else {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				}
			}
		}
	}
}
