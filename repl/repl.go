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
const CONT_PROMPT = "... "

// Start starts the interactive REPL session.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type RustX expressions or statements. Press Ctrl+C or Ctrl+D to exit.")

	var buffer strings.Builder
	braceDepth := 0
	parenDepth := 0

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

		// Count braces and parentheses outside strings
		braceDelta, parenDelta := countDelimiters(line)
		braceDepth += braceDelta
		parenDepth += parenDelta

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// If all open blocks/parens are closed, evaluate the buffer
		if braceDepth <= 0 && parenDepth <= 0 {
			input := buffer.String()
			buffer.Reset()
			braceDepth = 0
			parenDepth = 0

			trimmed := strings.TrimSpace(input)
			if trimmed == "" {
				continue
			}

			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func countDelimiters(line string) (int, int) {
	braceDelta := 0
	parenDelta := 0
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
		// Skip line comments
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			break
		}

		switch ch {
		case '{':
			braceDelta++
		case '}':
			braceDelta--
		case '(':
			parenDelta++
		case ')':
			parenDelta--
		}
	}

	return braceDelta, parenDelta
}

func printParserErrors(out io.Writer, syntaxErrors []*errors.SyntaxError) {
	for _, err := range syntaxErrors {
		fmt.Fprintln(out, err.Error())
	}
}
