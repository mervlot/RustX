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
	fmt.Fprintln(out, "Type 'exit()' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit()" {
			break
		}

		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		fullInput := buffer.String()

		if isInputIncomplete(fullInput) {
			fmt.Fprint(out, CONTINUATION_PROMPT)
			continue
		}

		buffer.Reset()

		l := lexer.New(fullInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func isInputIncomplete(input string) bool {
	openBraces := 0
	openParens := 0
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

	return openBraces > 0 || openParens > 0
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
