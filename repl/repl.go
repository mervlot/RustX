package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/token"
)

const PROMPT = ">>> "
const MULTI_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type some code or 'exit' to quit.\n")
	fmt.Fprintf(out, "%s", PROMPT)

	var buffer string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" && buffer == "" {
			break
		}

		if buffer == "" {
			buffer = line
		} else {
			buffer = buffer + "\n" + line
		}

		if !isBalanced(buffer) {
			fmt.Fprintf(out, "%s", MULTI_PROMPT)
			continue
		}

		// Buffer is balanced, evaluate it!
		l := lexer.New(buffer)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = ""
			fmt.Fprintf(out, "%s", PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}

		buffer = ""
		fmt.Fprintf(out, "%s", PROMPT)
	}
}

func isBalanced(input string) bool {
	l := lexer.New(input)
	braces := 0
	parens := 0
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}
		if tok.Type == token.LBRACE {
			braces++
		} else if tok.Type == token.RBRACE {
			braces--
		} else if tok.Type == token.LPAREN {
			parens++
		} else if tok.Type == token.RPAREN {
			parens--
		}
	}
	return braces <= 0 && parens <= 0
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
