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

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var lines []string
	openBraces := 0
	openParens := 0

	for {
		if openBraces == 0 && openParens == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if openBraces == 0 && openParens == 0 && strings.TrimSpace(line) == "exit" {
			return
		}

		lines = append(lines, line)

		// Count open/close braces and parentheses
		for _, ch := range line {
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

		// Continue reading if block or parens are still open
		if openBraces > 0 || openParens > 0 {
			continue
		}

		fullCode := strings.Join(lines, "\n")
		lines = []string{}

		if strings.TrimSpace(fullCode) == "" {
			continue
		}

		l := lexer.New(fullCode)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
