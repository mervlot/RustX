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
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
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

		input := buffer.String()

		if isCompleteInput(input) {
			l := lexer.New(input)
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

func isCompleteInput(input string) bool {
	var openBraces, openParens, openBrackets int
	inString := false
	var stringChar byte

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inString {
			if ch == '\\' {
				i++ // Skip escape
			} else if ch == stringChar {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			stringChar = ch
			continue
		}

		if ch == '/' && i+1 < len(input) {
			if input[i+1] == '/' {
				// Single line comment, skip till newline
				for i < len(input) && input[i] != '\n' {
					i++
				}
				continue
			}
			if input[i+1] == '*' {
				// Multiline comment, skip till */
				i += 2
				for i+1 < len(input) && !(input[i] == '*' && input[i+1] == '/') {
					i++
				}
				i++
				continue
			}
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
		case '[':
			openBrackets++
		case ']':
			openBrackets--
		}
	}

	return openBraces <= 0 && openParens <= 0 && openBrackets <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
