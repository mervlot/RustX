# RustX Starter

RustX Starter is a minimal, clean, and scalable interpreted programming language designed as a lightweight hybrid of Rust and Go. It combines the clean syntax and speed of Go with the expressive style and future extensibility of Rust.

---

## Architecture Overview

```
/token      - Token definitions, line/column location tracking, and keyword mappings.
/errors     - Human-readable SyntaxError and RuntimeError diagnostic types.
/lexer      - Lexical analyzer converting raw source code into tokens.
/ast        - Abstract Syntax Tree node definitions for expressions and statements.
/parser     - Pratt / Recursive Descent parser generating AST nodes with positional tracking.
/object     - Runtime value representations (Integer, Float, Boolean, String, Array, Function, Error, Builtin).
/evaluator  - Tree-walk evaluator with lexical scope environment management and standard library.
/repl       - Multiline interactive REPL prompt preserving state across commands.
/samples    - Sample programs demonstrating features, recursion, loops, and stdlib.
main.go     - CLI entry point supporting `run`, `repl`, `version`, and `build`.
```

---

## Language Features

- **Variables & Assignment**: Declarations with `var` and shorthand assignment with `:=`.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Relational (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), and Postfix (`++`, `--`).
- **Control Flow**: `if` / `else if` / `else`, `while` loops, and `for` loops (3-part headers and single condition headers).
- **Functions**: Named and anonymous functions, typed parameters, recursive call support, and optional main entry point auto-execution.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Standard Library**:
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `sleep`, `now`
  - Filesystem: `read_file`, `write_file`

---

## Getting Started & Build Instructions

### Prerequisites
- Go 1.20 or newer installed.

### Building
Compile the executable binary from the repository root:
```bash
go build -o rustx main.go
```

### Running Programs
Run a `.rx` script directly using the `run` command or shorthand:
```bash
./rustx run samples/hello.rx
# OR
./rustx samples/hello.rx
```

### Interactive REPL
Start the interactive shell:
```bash
./rustx
# OR
./rustx repl
```

Example session:
```text
RustX Starter v0.1
Type 'exit' or 'quit' to exit.
>>> x := 10
>>> func double(n int) int { return n * 2 }
>>> println(double(x))
20
```

### Checking Syntax (`build`)
Validate a script without running it:
```bash
./rustx build samples/factorial.rx
```

---

## Running Unit Tests

Run all unit tests across the codebase:
```bash
go test ./...
```
