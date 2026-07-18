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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type some code to execute or exit to quit.")

	var buffer []string
	openBraces := 0
	openParens := 0

	fmt.Fprint(out, PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if len(buffer) == 0 && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit") {
			break
		}

		buffer = append(buffer, line)

		// Count braces and parens for multiline tracking
		for _, char := range line {
			switch char {
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

		if openBraces <= 0 && openParens <= 0 {
			// All blocks and parens are balanced, so evaluate the buffer!
			code := strings.Join(buffer, "\n")
			buffer = nil // reset
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
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}

			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}
	}
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
