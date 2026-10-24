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
const INDENT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var multilineBuffer strings.Builder
	openBraces := 0
	openParens := 0

	fmt.Fprint(out, PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" || line == "quit" {
			break
		}

		multilineBuffer.WriteString(line + "\n")

		// Update bracket counts, ignoring parens/braces inside string literals
		openBraces, openParens = updateBracketCounts(line, openBraces, openParens)

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, INDENT_PROMPT)
			continue
		}

		fullInput := multilineBuffer.String()
		multilineBuffer.Reset()
		openBraces = 0
		openParens = 0

		if strings.TrimSpace(fullInput) == "" {
			fmt.Fprint(out, PROMPT)
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if errs := p.Errors(); len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(out, err.Error())
			}
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

func updateBracketCounts(line string, openBraces, openParens int) (int, int) {
	inString := false
	var stringQuote byte = 0

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inString {
			if ch == '\\' && i+1 < len(line) {
				i++ // skip escaped char
				continue
			}
			if ch == stringQuote {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			stringQuote = '"'
			continue
		}

		// Single line comment //
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			break
		}

		switch ch {
		case '{':
			openBraces++
		case '}':
			if openBraces > 0 {
				openBraces--
			}
		case '(':
			openParens++
		case ')':
			if openParens > 0 {
				openParens--
			}
		}
	}

	return openBraces, openParens
}
