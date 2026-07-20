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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type your code below. Semicolons are optional.")
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

			if len(p.Errors()) > 0 {
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
			fmt.Fprint(out, CONT_PROMPT)
		}
	}
}

func isBalanced(text string) bool {
	braces := 0
	parens := 0
	inString := false
	var prevChar rune

	for _, char := range text {
		if char == '"' && prevChar != '\\' {
			inString = !inString
		}
		if !inString {
			if char == '{' {
				braces++
			} else if char == '}' {
				braces--
			} else if char == '(' {
				parens++
			} else if char == ')' {
				parens--
			}
		}
		prevChar = char
	}

	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
