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
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()
		buffer = append(buffer, line)

		combined := strings.Join(buffer, "\n")
		if isBalanced(combined) {
			buffer = nil

			l := lexer.New(combined)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				fmt.Fprint(out, PROMPT)
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if err, ok := evaluated.(*object.Error); ok {
					fmt.Fprintln(out, err.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}
	}
}

func isBalanced(s string) bool {
	var braces int
	var parens int
	var brackets int
	inString := false
	var escape bool

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inString {
			if escape {
				escape = false
			} else if ch == '\\' {
				escape = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			continue
		}

		switch ch {
		case '{':
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		case '[':
			brackets++
		case ']':
			brackets--
		}
	}

	return braces <= 0 && parens <= 0 && brackets <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
