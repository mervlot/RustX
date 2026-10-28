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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line + "\n")

		// Count brace and paren balances
		openBraces += countCharUnescaped(line, '{') - countCharUnescaped(line, '}')
		openParens += countCharUnescaped(line, '(') - countCharUnescaped(line, ')')

		if openBraces <= 0 && openParens <= 0 {
			input := buffer.String()
			buffer.Reset()
			openBraces = 0
			openParens = 0

			if strings.TrimSpace(input) == "" {
				fmt.Fprint(out, PROMPT)
				continue
			}

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
				fmt.Fprintln(out, evaluated.Inspect())
			}
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func countCharUnescaped(s string, char byte) int {
	count := 0
	inString := false
	for i := 0; i < len(s); i++ {
		if s[i] == '"' && (i == 0 || s[i-1] != '\\') {
			inString = !inString
		}
		if !inString && s[i] == char {
			count++
		}
	}
	return count
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
