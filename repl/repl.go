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

	fmt.Fprintln(out, "RustX Starter REPL v0.1")
	fmt.Fprint(out, PROMPT)

	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		fullInput := strings.Join(lines, "\n")

		if !isBalanced(fullInput) {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		// Reset lines buffer
		lines = nil

		l := lexer.New(fullInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if errObj, ok := evaluated.(*object.Error); ok {
				fmt.Fprintf(out, "RuntimeError:\n%s\nat line %d column %d\n", errObj.Message, errObj.Token.Line, errObj.Token.Column)
			} else if evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}

func isBalanced(text string) bool {
	braces := 0
	parens := 0
	brackets := 0
	inString := false
	var escape bool

	for i := 0; i < len(text); i++ {
		ch := text[i]
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

		switch ch {
		case '"':
			inString = true
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
