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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")

	var inputBuffer strings.Builder
	openBraces := 0
	openParens := 0

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
		trimmed := strings.TrimSpace(line)

		if inputBuffer.Len() == 0 && (trimmed == "exit" || trimmed == "exit()") {
			fmt.Fprintln(out, "Goodbye!")
			return
		}

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		// Count brace and paren balances
		for _, ch := range line {
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

		// Keep taking input if blocks or parens are unclosed
		if openBraces > 0 || openParens > 0 {
			continue
		}

		fullInput := inputBuffer.String()
		inputBuffer.Reset()
		openBraces = 0
		openParens = 0

		if strings.TrimSpace(fullInput) == "" {
			continue
		}

		l := lexer.New(fullInput)
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

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
