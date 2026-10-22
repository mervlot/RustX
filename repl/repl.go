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

	fmt.Fprintln(out, "RustX Starter v0.1.0 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" || line == "quit" {
			break
		}

		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		if !isMultilineComplete(buffer.String()) {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		input := buffer.String()
		buffer.Reset()

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
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

		fmt.Fprint(out, PROMPT)
	}
}

func isMultilineComplete(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	var stringDelimiter byte

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inString {
			if ch == stringDelimiter && (i == 0 || input[i-1] != '\\') {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
			stringDelimiter = '"'
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

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(out, msg)
	}
}
