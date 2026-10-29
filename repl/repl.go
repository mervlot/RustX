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
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if buffer.Len() == 0 {
			trimmed := strings.TrimSpace(line)
			if trimmed == "exit" || trimmed == "quit" {
				fmt.Fprintln(out, "Goodbye!")
				break
			}
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		input := buffer.String()

		if isInputComplete(input) {
			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				for _, err := range p.Errors() {
					fmt.Fprintln(out, err.Error())
				}
				buffer.Reset()
				fmt.Fprint(out, PROMPT)
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			buffer.Reset()
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func isInputComplete(input string) bool {
	openBraces := 0
	openParens := 0
	openBrackets := 0
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
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			openBraces++
		case '}':
			openBraces--
		case '(':
			openParens++
		case ')':
			openParens--
		case '[':
			openBrackets++
		case ']':
			openBrackets--
		}
	}

	return openBraces <= 0 && openParens <= 0 && openBrackets <= 0 && !inString
}
