package repl

import (
	"bufio"
	"fmt"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"io"
	"strings"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		for !isBalanced(line) {
			fmt.Fprintf(out, "... ")
			scanned := scanner.Scan()
			if !scanned {
				break
			}
			line += "\n" + scanner.Text()
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func isBalanced(input string) bool {
	lbrace := strings.Count(input, "{")
	rbrace := strings.Count(input, "}")
	lparen := strings.Count(input, "(")
	rparen := strings.Count(input, ")")
	return lbrace == rbrace && lparen == rparen
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "SyntaxError:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
