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

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type RustX expressions or statements. Press Ctrl+C or Ctrl+D to exit.\n")
	fmt.Fprintf(out, "%s", PROMPT)

	var buffer strings.Builder
	braceCount := 0
	parenCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Count braces and parentheses outside of quotes
		inString := false
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
				continue
			}
			if inString {
				continue
			}
			switch ch {
			case '{':
				braceCount++
			case '}':
				braceCount--
			case '(':
				parenCount++
			case ')':
				parenCount--
			}
		}

		buffer.WriteString(line + "\n")

		// If multiline statement is open, continue prompting
		if braceCount > 0 || parenCount > 0 {
			fmt.Fprintf(out, "%s", CONT_PROMPT)
			continue
		}

		fullInput := buffer.String()
		buffer.Reset()
		braceCount = 0
		parenCount = 0

		trimmed := strings.TrimSpace(fullInput)
		if trimmed == "" {
			fmt.Fprintf(out, "%s", PROMPT)
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintf(out, "%s\n", err.Error())
			}
			fmt.Fprintf(out, "%s", PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}

		fmt.Fprintf(out, "%s", PROMPT)
	}
}
