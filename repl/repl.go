package repl

import (
	"bufio"
	"fmt"
	"io"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const CONT_PROMPT = "... "

// Start launches the interactive REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' to quit.")

	var buffer strings.Builder
	inMultiLine := false

	for {
		if !inMultiLine {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if !inMultiLine && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit") {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		currentCode := buffer.String()

		// Check if input is complete (balanced braces and parentheses)
		if isCompleteInput(currentCode) {
			inMultiLine = false
			buffer.Reset()

			l := lexer.New(currentCode)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				for _, err := range p.Errors() {
					fmt.Fprintln(out, err.Error())
				}
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
		} else {
			inMultiLine = true
		}
	}
}

func isCompleteInput(code string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(code); i++ {
		ch := code[i]
		if ch == '"' && (i == 0 || code[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		switch ch {
		case '{':
			openBraces++
		case '}':
			if openBraces > 0 {
				openBraces--
			}
		case '(':
			openParens++
		case ')':
			if openParens > 0 {
				openParens--
			}
		}
	}

	return openBraces == 0 && openParens == 0
}
