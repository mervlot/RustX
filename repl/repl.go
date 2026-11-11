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

const (
	PROMPT     = ">>> "
	PROMPT_SUB = "... "
)

// Start starts an interactive REPL reading from in and writing to out.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type expressions or statements. Press Ctrl+D (or type 'exit') to quit.\n")

	var buffer strings.Builder
	depth := 0

	for {
		if depth == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, PROMPT_SUB)
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

		// Update bracket and paren depth
		depth += countBalance(line)
		if depth < 0 {
			depth = 0
		}

		if depth > 0 {
			continue
		}

		code := buffer.String()
		buffer.Reset()

		if strings.TrimSpace(code) == "" {
			continue
		}

		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintln(out, err.Error())
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

// countBalance computes net open braces/parentheses, ignoring strings.
func countBalance(line string) int {
	balance := 0
	inString := false
	var prev byte

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '"' && prev != '\\' {
			inString = !inString
		} else if !inString {
			// check for single line comment start
			if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				break
			}
			switch ch {
			case '{', '(':
				balance++
			case '}', ')':
				balance--
			}
		}
		prev = ch
	}

	return balance
}
