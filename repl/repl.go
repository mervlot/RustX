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

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var inputBuffer strings.Builder

	for {
		if inputBuffer.Len() == 0 {
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

		if inputBuffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			fmt.Fprintln(out, "Goodbye!")
			return
		}

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		bufferStr := inputBuffer.String()
		if isInputComplete(bufferStr) {
			l := lexer.New(bufferStr)
			p, err := parser.New(l)
			if err != nil {
				fmt.Fprintln(out, err.Error())
				inputBuffer.Reset()
				continue
			}

			program, err := p.ParseProgram()
			if err != nil {
				fmt.Fprintln(out, err.Error())
				inputBuffer.Reset()
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

			inputBuffer.Reset()
		}
	}
}

func isInputComplete(code string) bool {
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
			openBraces--
		case '(':
			openParens++
		case ')':
			openParens--
		}
	}

	return openBraces <= 0 && openParens <= 0
}
