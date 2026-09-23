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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or Ctrl+C to quit.")

	var buffer strings.Builder
	inMultiline := false

	for {
		if inMultiline {
			fmt.Fprint(out, CONT_PROMPT)
		} else {
			fmt.Fprint(out, PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if !inMultiline && strings.TrimSpace(line) == "exit" {
			return
		}

		if inMultiline {
			buffer.WriteString("\n")
		} else {
			buffer.Reset()
		}
		buffer.WriteString(line)

		code := buffer.String()
		if isCompleteCode(code) {
			inMultiline = false
			buffer.Reset()

			l := lexer.New(code)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			}
		} else {
			inMultiline = true
		}
	}
}

func isCompleteCode(code string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(code); i++ {
		ch := code[i]
		if ch == '"' && (i == 0 || code[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
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

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
