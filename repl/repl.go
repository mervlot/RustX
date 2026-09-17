package repl

import (
	"bufio"
	"fmt"
	"io"

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
	fmt.Fprintln(out, "Type 'exit' or 'quit' to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer string
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if buffer == "" && (line == "exit" || line == "quit") {
			break
		}

		buffer += line + "\n"

		// Count open braces and parens (naive count ignoring strings/comments)
		openBraces += countChar(line, '{') - countChar(line, '}')
		openParens += countChar(line, '(') - countChar(line, ')')

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		l := lexer.New(buffer)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = ""
			openBraces = 0
			openParens = 0
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		buffer = ""
		openBraces = 0
		openParens = 0
		fmt.Fprint(out, PROMPT)
	}
}

func countChar(s string, char byte) int {
	count := 0
	inString := false
	for i := 0; i < len(s); i++ {
		if s[i] == '"' && (i == 0 || s[i-1] != '\\') {
			inString = !inString
		}
		if !inString && s[i] == char {
			count++
		}
	}
	return count
}

func printParserErrors(out io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
