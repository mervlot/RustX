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
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter REPL v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Check if braces / parens are balanced
		if !isBalanced(buffer.String()) {
			fmt.Fprint(out, CONTINUATION_PROMPT)
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
			printParserErrors(out, p)
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

func isBalanced(code string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	escaped := false

	for i := 0; i < len(code); i++ {
		ch := code[i]
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
			openBraces--
		case '(':
			openParens++
		case ')':
			openParens--
		}
	}

	return openBraces <= 0 && openParens <= 0 && !inString
}

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err.Error())
	}
}
