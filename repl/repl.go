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

	fmt.Fprintln(out, "RustX Starter REPL v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if buffer.Len() == 0 && strings.TrimSpace(line) == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		bufStr := buffer.String()
		if isUnmatched(bufStr) {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		l := lexer.New(bufStr)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p)
			buffer.Reset()
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		buffer.Reset()
		fmt.Fprint(out, PROMPT)
	}
}

func isUnmatched(input string) bool {
	openBraces := 0
	openParens := 0
	openBrackets := 0
	inString := false
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
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
		case '[':
			openBrackets++
		case ']':
			if openBrackets > 0 {
				openBrackets--
			}
		}
	}

	return openBraces > 0 || openParens > 0 || openBrackets > 0 || inString
}

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err.Error())
	}
}
