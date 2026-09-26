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

	fmt.Fprintln(out, "RustX Starter v0.1.0 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if buffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		parenCount, braceCount := countDelimiters(buffer.String())

		if parenCount > 0 || braceCount > 0 {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		code := buffer.String()
		buffer.Reset()

		if strings.TrimSpace(code) == "" {
			fmt.Fprint(out, PROMPT)
			continue
		}

		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Fprintln(out, err.Error())
			}
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func countDelimiters(code string) (int, int) {
	parens := 0
	braces := 0
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
		case '(':
			parens++
		case ')':
			if parens > 0 {
				parens--
			}
		case '{':
			braces++
		case '}':
			if braces > 0 {
				braces--
			}
		}
	}

	return parens, braces
}
