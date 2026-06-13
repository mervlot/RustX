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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, "... ")
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		buffer.WriteString(line + "\n")

		if shouldEvaluate(buffer.String()) {
			l := lexer.New(buffer.String())
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer.Reset()
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
			buffer.Reset()
		}
	}
}

func shouldEvaluate(input string) bool {
	// Simple heuristic: check for balanced braces
	l := lexer.New(input)
	lbrace := 0
	rbrace := 0
	lparen := 0
	rparen := 0
	for {
		tok := l.NextToken()
		if tok.Type == "EOF" {
			break
		}
		if tok.Type == "{" {
			lbrace++
		}
		if tok.Type == "}" {
			rbrace++
		}
		if tok.Type == "(" {
			lparen++
		}
		if tok.Type == ")" {
			rparen++
		}
	}
	return lbrace <= rbrace && lparen <= rparen
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "SyntaxError:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
