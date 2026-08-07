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
const MULTILINE_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	ev := evaluator.New("repl")

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type expressions or statements. Press Enter twice to submit multiline blocks.\n")

	var lines []string

	for {
		if len(lines) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTILINE_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// If we are building a block and the user hits Enter on an empty line, evaluate what we have.
		if len(lines) > 0 && strings.TrimSpace(line) == "" {
			fullInput := strings.Join(lines, "\n")
			lines = nil
			evaluate(fullInput, ev, env, out)
			continue
		}

		lines = append(lines, line)

		// Check if parentheses and curly braces are completely balanced.
		// If balanced and we are not in the middle of a block, evaluate.
		fullInput := strings.Join(lines, "\n")
		if isBalanced(fullInput) {
			lines = nil
			evaluate(fullInput, ev, env, out)
		}
	}
}

func isBalanced(input string) bool {
	// A simple helper to balance brackets to determine if we should wait for more input or evaluate immediately
	braces := 0
	parens := 0
	inString := false
	var stringChar rune

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if inString {
			if r == '\\' && i+1 < len(runes) {
				i++ // skip escaped char
			} else if r == stringChar {
				inString = false
			}
		} else {
			switch r {
			case '"', '\'':
				inString = true
				stringChar = r
			case '{':
				braces++
			case '}':
				braces--
			case '(':
				parens++
			case ')':
				parens--
			}
		}
	}

	return braces <= 0 && parens <= 0
}

func evaluate(input string, ev *evaluator.Evaluator, env *object.Environment, out io.Writer) {
	if strings.TrimSpace(input) == "" {
		return
	}

	l := lexer.New(input)
	p := parser.New(l, "repl")

	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := ev.Eval(program, env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		if errObj, ok := evaluated.(*object.Error); ok {
			fmt.Fprintf(out, "%s\n", errObj.Inspect())
		} else {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
