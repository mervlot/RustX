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
	fmt.Fprintln(out, "Type 'exit' or Ctrl+C to exit.")

	var buffer strings.Builder
	openBraces := 0
	openParens := 0

	for {
		if buffer.Len() == 0 {
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

		if buffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		buffer.WriteString(line + "\n")

		// Count brace and paren balances while ignoring contents inside strings
		openBraces, openParens = countBalances(buffer.String())

		if openBraces > 0 || openParens > 0 {
			continue
		}

		input := buffer.String()
		buffer.Reset()

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else if evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func countBalances(input string) (braces int, parens int) {
	inString := false
	var stringChar byte

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '\\' {
				i++ // skip escaped char
			} else if ch == stringChar {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			stringChar = ch
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

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err.Error())
	}
}
