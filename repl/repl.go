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

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")
	fmt.Fprint(out, PROMPT)

	var inputBuffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" && inputBuffer.Len() == 0 {
			break
		}

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		bufferStr := inputBuffer.String()

		// Check if braces/parens are balanced for multiline input
		if !isBalanced(bufferStr) {
			fmt.Fprint(out, CONTINUATION_PROMPT)
			continue
		}

		l := lexer.New(bufferStr)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintln(out, err.Error())
			}
			inputBuffer.Reset()
			fmt.Fprint(out, PROMPT)
			continue
		}

		// Intercept stdout calls in REPL if needed or let builtins write to stdout/out
		evaluated := evaluator.EvalWithOutput(program, env, out)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		inputBuffer.Reset()
		fmt.Fprint(out, PROMPT)
	}
}

func isBalanced(code string) bool {
	openBraces := 0
	openParens := 0
	openBrackets := 0
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
		case '[':
			openBrackets++
		case ']':
			openBrackets--
		}
	}

	return openBraces <= 0 && openParens <= 0 && openBrackets <= 0
}
