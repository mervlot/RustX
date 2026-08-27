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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			break
		}

		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		bufStr := buffer.String()
		if isCompleteBlock(bufStr) {
			l := lexer.New(bufStr)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				buffer.Reset()
				fmt.Fprint(out, PROMPT)
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				fmt.Fprintln(out, evaluated.Inspect())
			}

			buffer.Reset()
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func isCompleteBlock(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	var stringDelimiter byte

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inString {
			if ch == stringDelimiter && (i == 0 || input[i-1] != '\\') {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			stringDelimiter = '"'
			continue
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

	return openBraces <= 0 && openParens <= 0 && !inString
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
