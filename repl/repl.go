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

	fmt.Fprintln(out, "RustX Starter REPL v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")

	var buffer strings.Builder
	inMultiLine := false

	for {
		if !inMultiLine {
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

		if !inMultiLine && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		bufStr := buffer.String()
		if isBalanced(bufStr) {
			l := lexer.New(bufStr)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer.Reset()
				inMultiLine = false
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

			buffer.Reset()
			inMultiLine = false
		} else {
			inMultiLine = true
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}

func isBalanced(code string) bool {
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
