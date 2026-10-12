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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Count depth of braces and parens to support multiline inputs
		inString := false
		for _, ch := range line {
			if ch == '"' {
				inString = !inString
			}
			if !inString {
				if ch == '{' {
					openBraces++
				} else if ch == '}' && openBraces > 0 {
					openBraces--
				} else if ch == '(' {
					openParens++
				} else if ch == ')' && openParens > 0 {
					openParens--
				}
			}
		}

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		code := buffer.String()
		buffer.Reset()
		openBraces = 0
		openParens = 0

		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintf(out, "RuntimeError: %s\n", evaluated.Inspect())
			} else if evaluated.Type() != object.RETURN_VALUE_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
