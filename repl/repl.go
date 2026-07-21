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
	fmt.Fprintln(out, "Type some code to evaluate or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()
		buffer = append(buffer, line)

		combined := strings.Join(buffer, "\n")
		if isBalanced(combined) {
			// Ready to evaluate
			l := lexer.New(combined)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			buffer = nil
			fmt.Fprint(out, PROMPT)
		} else {
			// Ask for more input
			fmt.Fprint(out, MULTI_PROMPT)
		}
	}
}

func isBalanced(text string) bool {
	braces := 0
	parens := 0
	inString := false
	var escape bool

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if escape {
			escape = false
			continue
		}
		if ch == '\\' && inString {
			escape = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if !inString {
			if ch == '{' {
				braces++
			} else if ch == '}' {
				braces--
			} else if ch == '(' {
				parens++
			} else if ch == ')' {
				parens--
			}
		}
	}
	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%v\n", err)
	}
}
