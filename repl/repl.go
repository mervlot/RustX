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

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
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
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintln(out, errObj.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
			buffer.Reset()
		}
	}
}

func isBalanced(input string) bool {
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

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
