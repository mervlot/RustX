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

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type 'exit' to quit.\n\n")

	var buffer []string

	for {
		if len(buffer) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// Allow quitting from REPL
		if len(buffer) == 0 && strings.TrimSpace(line) == "exit" {
			return
		}

		buffer = append(buffer, line)
		fullInput := strings.Join(buffer, "\n")

		if isBalanced(fullInput) {
			// Ready to parse and execute!
			l := lexer.New(fullInput)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p)
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
					fmt.Fprintf(out, "%s\n", evaluated.Inspect())
				}
			}

			// Clear the buffer
			buffer = nil
		}
	}
}

func isBalanced(text string) bool {
	lbrace := 0
	lparen := 0
	lbracket := 0
	inString := false

	for i := 0; i < len(text); i++ {
		c := text[i]
		if c == '"' {
			if i > 0 && text[i-1] == '\\' {
				continue
			}
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		switch c {
		case '{':
			lbrace++
		case '}':
			lbrace--
		case '(':
			lparen++
		case ')':
			lparen--
		case '[':
			lbracket++
		case ']':
			lbracket--
		}
	}
	return lbrace <= 0 && lparen <= 0 && lbracket <= 0
}

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
