package repl

import (
	"bufio"
	"fmt"
	"io"

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

	var buffer string
	openBraces := 0
	openParens := 0

	fmt.Fprint(out, PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if buffer == "" && (line == "exit" || line == "quit") {
			break
		}

		buffer += line + "\n"

		// Calculate open braces and parens (ignore inside string literals)
		inString := false
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
				continue
			}
			if !inString {
				if ch == '{' {
					openBraces++
				} else if ch == '}' && openBraces > 0 {
					openBraces--
				} else if ch == '(' {
					openParens++
				} else if ch == ')' && openParens > 0 {
					openParens--
				}
			}
		}

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		l := lexer.New(buffer)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = ""
			openBraces = 0
			openParens = 0
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		buffer = ""
		openBraces = 0
		openParens = 0
		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
