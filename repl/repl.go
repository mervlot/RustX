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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1.0 Interactive REPL\n")
	fmt.Fprintf(out, "Type statements or expressions. Press Ctrl+C or Ctrl+D to exit.\n\n")

	var buffer strings.Builder

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
		buffer.WriteString(line + "\n")

		bufStr := buffer.String()
		if isInputComplete(bufStr) {
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
				fmt.Fprintln(out, evaluated.Inspect())
			}

			buffer.Reset()
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

	return openBraces <= 0 && openParens <= 0 && openBrackets <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
