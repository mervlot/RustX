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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if buffer.Len() == 0 {
			trimmed := strings.TrimSpace(line)
			if trimmed == "exit" || trimmed == "quit" {
				break
			}
		}

		buffer.WriteString(line + "\n")
		text := buffer.String()

		if isMultilineBlock(text) {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		l := lexer.New(text)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				fmt.Fprintln(out, msg)
			}
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

func isMultilineBlock(code string) bool {
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

	return openBraces > 0 || openParens > 0
}
