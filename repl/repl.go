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
	fmt.Fprintln(out, "Type \"exit\" or press Ctrl+C to quit.")

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
		trimmed := strings.TrimSpace(line)
		if buffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		fullInput := buffer.String()

		if isBalanced(fullInput) {
			l := lexer.New(fullInput)
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
					fmt.Fprintln(out, evaluated.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			buffer.Reset()
		}
	}
}

func isBalanced(input string) bool {
	braceCount := 0
	parenCount := 0
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
			braceCount++
		case '}':
			braceCount--
		case '(':
			parenCount++
		case ')':
			parenCount--
		}
	}

	return braceCount <= 0 && parenCount <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
