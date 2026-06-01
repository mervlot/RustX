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
const CONTINUE_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	var inputBuffer strings.Builder

	for {
		if inputBuffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUE_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		combinedInput := inputBuffer.String()

		if !isBalanced(combinedInput) {
			continue
		}

		l := lexer.New(combinedInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			// If it's a simple syntax error because of incomplete input, we might want to continue.
			// But isBalanced should handle most cases.
			printParserErrors(out, p.Errors())
			inputBuffer.Reset()
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
		inputBuffer.Reset()
	}
}

func isBalanced(input string) bool {
	l := lexer.New(input)
	lbraceCount := 0
	lparenCount := 0

	for {
		tok := l.NextToken()
		if tok.Type == "EOF" {
			break
		}
		switch tok.Type {
		case "{":
			lbraceCount++
		case "}":
			lbraceCount--
		case "(":
			lparenCount++
		case ")":
			lparenCount--
		}
	}

	return lbraceCount <= 0 && lparenCount <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
