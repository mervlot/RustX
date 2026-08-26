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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Count brace and paren balance, skipping characters in string literals
		inString := false
		escaped := false
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' && inString {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = !inString
				continue
			}
			if !inString {
				switch ch {
				case '{':
					openBraces++
				case '}':
					openBraces--
				case '(':
					openParens++
				case ')':
					openParens--
				}
			}
		}

		if openBraces <= 0 && openParens <= 0 {
			input := buffer.String()
			buffer.Reset()
			openBraces = 0
			openParens = 0

			l := lexer.New(input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
