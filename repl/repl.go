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
	fmt.Fprintln(out, "Type commands or definitions. Press Enter to evaluate.")

	var buffer []string
	braceBalance := 0
	parenBalance := 0

	for {
		if len(buffer) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		buffer = append(buffer, line)

		b, p := countBalance(line)
		braceBalance += b
		parenBalance += p

		if braceBalance <= 0 && parenBalance <= 0 {
			code := strings.Join(buffer, "\n")
			buffer = nil
			braceBalance = 0
			parenBalance = 0

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
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintln(out, errObj.Message)
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
		}
	}
}

func countBalance(line string) (int, int) {
	braces := 0
	parens := 0
	inString := false
	escaped := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if !inString {
			// Ignore trailing single-line comments in counting braces/parens
			if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				break
			}
			if ch == '{' {
				braces++
			} else if ch == '}' {
				braces--
			} else if ch == '(' {
				parens++
			} else if ch == ')' {
				parens--
			}
		}
	}
	return braces, parens
}
