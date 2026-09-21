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
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		if isInputComplete(buffer.String()) {
			code := buffer.String()
			buffer.Reset()

			l := lexer.New(code)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				fmt.Fprint(out, PROMPT)
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
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func isInputComplete(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '"' && (i == 0 || input[i-1] != '\\') {
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
			openBraces--
		case '(':
			openParens++
		case ')':
			openParens--
		}
	}

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
