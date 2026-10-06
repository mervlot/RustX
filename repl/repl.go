package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/username/rustx/errors"
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

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		if isBalanced(buffer.String()) {
			input := buffer.String()
			buffer.Reset()

			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintln(out, errObj.Inspect())
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
		}
	}
}

func isBalanced(input string) bool {
	braces := 0
	parens := 0
	inString := false
	var prevChar rune

	for _, ch := range input {
		if ch == '"' && prevChar != '\\' {
			inString = !inString
		}

		if !inString {
			switch ch {
			case '{':
				braces++
			case '}':
				braces--
			case '(':
				parens++
			case ')':
				parens--
			}
		}
		prevChar = ch
	}

	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
