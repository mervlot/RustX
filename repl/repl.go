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

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Count open braces / parens to determine if multiline input is incomplete
		openBraces += countCharUnescaped(line, '{') - countCharUnescaped(line, '}')
		openParens += countCharUnescaped(line, '(') - countCharUnescaped(line, ')')

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, CONTINUATION_PROMPT)
			continue
		}

		input := buffer.String()
		buffer.Reset()
		openBraces = 0
		openParens = 0

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

func printParserErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
