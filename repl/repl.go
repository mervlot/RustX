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
const CONTINUE_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type commands or expressions. Press Ctrl+C to exit.\n")

	for {
		fmt.Fprintf(out, PROMPT)
		var lines []string

		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)

			fullInput := strings.Join(lines, "\n")
			if isComplete(fullInput) {
				break
			}
			fmt.Fprintf(out, CONTINUE_PROMPT)
		}

		input := strings.Join(lines, "\n")
		if input == "" {
			continue
		}

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func isComplete(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '"' {
			if i > 0 && input[i-1] != '\\' {
				inString = !inString
			} else if i == 0 {
				inString = !inString
			}
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
	}

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg + "\n")
	}
}
