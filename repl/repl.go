package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/username/rustx/errors"
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

	fmt.Fprintln(out, "RustX Starter v0.1 Prompt")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var inputBuffer strings.Builder

	for {
		if inputBuffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if inputBuffer.Len() == 0 && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit") {
			return
		}

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		bufferStr := inputBuffer.String()

		if isInputComplete(bufferStr) {
			l := lexer.New(bufferStr)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				inputBuffer.Reset()
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				} else if rv, ok := evaluated.(*object.ReturnValue); ok {
					fmt.Fprintln(out, rv.Value.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			inputBuffer.Reset()
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

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
