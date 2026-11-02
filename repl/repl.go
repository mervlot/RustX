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
const CONTINUE_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or Ctrl+D to quit.")

	var inputBuffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if inputBuffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUE_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if inputBuffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		// Count brace balance (ignoring braces within strings)
		deltaBraces, deltaParens := countDelimiters(line)
		openBraces += deltaBraces
		openParens += deltaParens

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		if openBraces <= 0 && openParens <= 0 {
			fullInput := inputBuffer.String()
			inputBuffer.Reset()
			openBraces = 0
			openParens = 0

			if strings.TrimSpace(fullInput) == "" {
				continue
			}

			l := lexer.New(fullInput)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				} else if evaluated.Type() != object.NULL_OBJ {
					// In REPL, if an expression evaluated to non-null, print it
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
		}
	}
}

func countDelimiters(line string) (int, int) {
	deltaBraces := 0
	deltaParens := 0
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
			deltaBraces++
		case '}':
			deltaBraces--
		case '(':
			deltaParens++
		case ')':
			deltaParens--
		}
	}
	return deltaBraces, deltaParens
}

func printParserErrors(out io.Writer, errors []*errors.SyntaxError) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
