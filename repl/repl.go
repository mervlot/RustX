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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

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
		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
			break
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
				if evaluated.Type() == object.ERROR_OBJ {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				} else {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				}
			}

			buffer.Reset()
		}
	}
}

func isBalanced(code string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	var prevChar byte

	for i := 0; i < len(code); i++ {
		ch := code[i]
		if ch == '"' && prevChar != '\\' {
			inString = !inString
		}
		if !inString {
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
		prevChar = ch
	}

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
