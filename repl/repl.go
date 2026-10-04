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
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if openBraces == 0 && openParens == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if openBraces == 0 && openParens == 0 && (trimmed == "exit" || trimmed == "quit") {
			fmt.Fprintln(out, "Goodbye!")
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Track block balancing (ignoring inside string literals)
		ob, op := countUnmatchedBrackets(buffer.String())
		openBraces = ob
		openParens = op

		if openBraces > 0 || openParens > 0 {
			continue
		}

		input := buffer.String()
		buffer.Reset()

		if strings.TrimSpace(input) == "" {
			continue
		}

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintf(out, "%s\n", err.Error())
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}
	}
}

func countUnmatchedBrackets(input string) (int, int) {
	braces := 0
	parens := 0
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
