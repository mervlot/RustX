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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type RustX expressions or statements. Press Ctrl+D (or Ctrl+C) to exit.")

	var buffer strings.Builder
	braceCount := 0
	parenCount := 0

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// Count balance of braces and parens ignoring strings
		braceDelta, parenDelta := countDelimiters(line)
		braceCount += braceDelta
		parenCount += parenDelta

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// If all brackets and parens are closed, evaluate the accumulated code
		if braceCount <= 0 && parenCount <= 0 {
			input := buffer.String()
			buffer.Reset()
			braceCount = 0
			parenCount = 0

			trimmed := strings.TrimSpace(input)
			if trimmed == "" {
				continue
			}

			l := lexer.New(input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				for _, err := range p.Errors() {
					fmt.Fprintln(out, err.Error())
				}
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func countDelimiters(s string) (int, int) {
	braces := 0
	parens := 0
	inString := false
	var prev byte

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' && prev != '\\' {
			inString = !inString
		}
		if !inString {
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
		prev = ch
	}

	return braces, parens
}
