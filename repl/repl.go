package repl

import (
	"bufio"
	"fmt"
	"io"
	rxErrors "rustx/errors"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const CONTINUE_PROMPT = "... "

// Start starts the interactive REPL session.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C/Ctrl+D to exit.")

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUE_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if buffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		if !isInputComplete(buffer.String()) {
			continue
		}

		fullInput := buffer.String()
		buffer.Reset()

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func isInputComplete(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	var prevChar rune

	for _, ch := range input {
		if ch == '"' && prevChar != '\\' {
			inString = !inString
		}
		if !inString {
			switch ch {
			case '{':
				openBraces++
			case '}':
				openBraces--
			case '(':
				openParens++
			case ')':
				openParens--
			}
		}
		prevChar = ch
	}

	return openBraces <= 0 && openParens <= 0 && !inString
}

func printParserErrors(out io.Writer, errors []*rxErrors.SyntaxError) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg.Error())
	}
}
