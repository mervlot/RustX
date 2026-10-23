package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
)

const PROMPT = ">>> "
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter REPL v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+D to quit.")

	var lines []string

	for {
		if len(lines) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if len(lines) == 0 && trimmed == "exit" {
			return
		}

		lines = append(lines, line)
		fullInput := strings.Join(lines, "\n")

		if isBlockClosed(fullInput) {
			l := lexer.New(fullInput)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				lines = nil
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}

			lines = nil
		}
	}
}

func isBlockClosed(input string) bool {
	braces := 0
	parens := 0
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
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		}
	}

	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(out, msg)
	}
}
