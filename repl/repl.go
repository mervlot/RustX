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
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Track braces & parens state
		openBraces, openParens = countUnbalancedDelimiters(line, openBraces, openParens)

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, CONTINUATION_PROMPT)
			continue
		}

		code := buffer.String()
		buffer.Reset()
		openBraces = 0
		openParens = 0

		if strings.TrimSpace(code) == "" {
			fmt.Fprint(out, PROMPT)
			continue
		}

		l := lexer.New(code)
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

func countUnbalancedDelimiters(line string, initialBraces, initialParens int) (int, int) {
	braces := initialBraces
	parens := initialParens
	inString := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '"' && (i == 0 || line[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		if ch == '{' {
			braces++
		} else if ch == '}' {
			if braces > 0 {
				braces--
			}
		} else if ch == '(' {
			parens++
		} else if ch == ')' {
			if parens > 0 {
				parens--
			}
		}
	}

	return braces, parens
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
