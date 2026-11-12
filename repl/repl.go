package repl

import (
	"bufio"
	"fmt"
	"io"
	"rustx/errors"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const CONTINUATION_PROMPT = "... "

// Start starts the interactive REPL session.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type RustX expressions or statements. Press Ctrl+D or type exit to quit.\n")

	var buffer strings.Builder
	depth := 0

	for {
		if depth == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if depth == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Calculate brace depth (ignoring strings)
		depth = computeBraceDepth(buffer.String())
		if depth > 0 {
			continue
		}

		input := buffer.String()
		buffer.Reset()
		depth = 0

		if strings.TrimSpace(input) == "" {
			continue
		}

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else if evaluated != evaluator.NULL {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func computeBraceDepth(input string) int {
	depth := 0
	inString := false
	var prevChar rune

	for _, ch := range input {
		if ch == '"' && prevChar != '\\' {
			inString = !inString
		}
		if !inString {
			if ch == '{' || ch == '(' {
				depth++
			} else if ch == '}' || ch == ')' {
				depth--
			}
		}
		prevChar = ch
	}

	if depth < 0 {
		return 0
	}
	return depth
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
