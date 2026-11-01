package repl

import (
	"bufio"
	"fmt"
	"io"
	"rustx/errors"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const PROMPT_CONT = "... "

// Start starts the interactive REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type statements or expressions. Press Ctrl+C or Ctrl+D to exit.")

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, PROMPT_CONT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit" {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		currentText := buffer.String()
		if isInputComplete(currentText) {
			input := strings.TrimSpace(currentText)
			buffer.Reset()

			if input == "" {
				continue
			}

			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				} else if evaluated.Type() != object.NULL_OBJ {
					// In REPL, print non-null expression result
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
		}
	}
}

func isInputComplete(s string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' && (i == 0 || s[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if ch == '{' {
			openBraces++
		} else if ch == '}' {
			openBraces--
		} else if ch == '(' {
			openParens++
		} else if ch == ')' {
			openParens--
		}
	}

	return openBraces <= 0 && openParens <= 0 && !inString
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
