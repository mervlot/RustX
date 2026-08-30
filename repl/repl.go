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
	fmt.Fprintln(out, "Type statements or expressions to evaluate. Multiline blocks supported.")

	var inputBuffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if inputBuffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		// Track nested braces and parens for multiline input
		openBraces, openParens = countUnmatched(inputBuffer.String())

		if openBraces > 0 || openParens > 0 {
			continue
		}

		code := inputBuffer.String()
		inputBuffer.Reset()

		if strings.TrimSpace(code) == "" {
			continue
		}

		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}

func countUnmatched(s string) (int, int) {
	braces := 0
	parens := 0
	inString := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' && (i == 0 || s[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		switch ch {
		case '{':
			braces++
		case '}':
			if braces > 0 {
				braces--
			}
		case '(':
			parens++
		case ')':
			if parens > 0 {
				parens--
			}
		}
	}

	return braces, parens
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(out, msg)
	}
}
