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
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	inMultiLine := false

	for scanner.Scan() {
		line := scanner.Text()

		if !inMultiLine && strings.TrimSpace(line) == "exit" {
			break
		}

		if inMultiLine {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		fullInput := buffer.String()

		if isBalanced(fullInput) {
			l := lexer.New(fullInput)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated != evaluator.NULL {
					if errObj, ok := evaluated.(*object.Error); ok {
						fmt.Fprintf(out, "RuntimeError: %s\n", errObj.Message)
					} else {
						fmt.Fprintln(out, evaluated.Inspect())
					}
				}
			}

			buffer.Reset()
			inMultiLine = false
			fmt.Fprint(out, PROMPT)
		} else {
			inMultiLine = true
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func isBalanced(code string) bool {
	openBraces := 0
	openParens := 0
	openBrackets := 0
	inString := false
	var stringDelimiter byte = 0
	escaped := false

	for i := 0; i < len(code); i++ {
		ch := code[i]

		if inString {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == stringDelimiter {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
			stringDelimiter = ch
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
