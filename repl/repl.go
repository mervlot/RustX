package repl

import (
	"bufio"
	"fmt"
	"io"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const MULTI_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit the REPL.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	inMultiLine := false

	for scanner.Scan() {
		line := scanner.Text()

		if !inMultiLine && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit") {
			break
		}

		buffer.WriteString(line + "\n")
		bufStr := buffer.String()

		if isMultilinePending(bufStr) {
			inMultiLine = true
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		inMultiLine = false
		l := lexer.New(bufStr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p)
			buffer.Reset()
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

		buffer.Reset()
		fmt.Fprint(out, PROMPT)
	}
}

func isMultilinePending(code string) bool {
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

	return openBraces > 0 || openParens > 0
}

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err.Error())
	}
}
