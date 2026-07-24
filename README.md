# RustX Starter

Welcome to **RustX Starter**, a minimal but clean, scalable programming language inspired by Rust and Go. It is designed to be elegant, fast, and readable, focusing on modern language constructs as a lightweight hybrid of scripting and systems programming.

This implementation acts as the foundation-level architecture, designed so that advanced systems (like ownership, borrow checking, and lifetimes) can be integrated in future phases without a full rewrite of the lexer, parser, or interpreter.

---

## Architecture Overview

The system is implemented in Go and follows a highly modular, decoupled package architecture:

*   **`token/`**: Defines lexical token types (e.g., keywords, operators, literals, delimiters).
*   **`errors/`**: Custom diagnostic structures (`SyntaxError`, `RuntimeError`) for friendly, position-aware error reporting with line/column information.
*   **`lexer/`**: Lexical analyzer that breaks source text into tokens, skipping single-line (`//`) and multi-line (`/* */`) comments, and processing string escape sequences (e.g., `\n`, `\t`).
*   **`ast/`**: Extensible Abstract Syntax Tree nodes that all implement position reporting (`GetToken()`) for precise error locations.
*   **`parser/`**: Pratt recursive descent parser handling optional semicolons, expression priorities, prefix/infix/postfix operations, and unparenthesized loop headers/conditionals.
*   **`object/`**: The runtime-represented objects (e.g., numbers, strings, functions, arrays) and scoped variable Environments.
*   **`evaluator/`**: An AST-walking tree interpreter with short-circuiting logical operations, division-by-zero prevention, and rich built-in functions.
*   **`repl/`**: A multiline interactive shell that balances open braces/parentheses and preserves state.
*   **`samples/`**: Example `.rx` source files demonstrating language features.

---

## Language Specifications

### Syntax & Types
*   **Keywords**: `func`, `var`, `if`, `else`, `return`, `while`, `for`, `true`, `false`.
*   **Variables**: Mutable by default, defined via `var x = 5` or shorthand `x := 5`.
*   **Types**: `int` (Integer), `float` (Float), `bool` (Boolean), `string` (String), and `Array` (via `[elements]`).
*   **Operators**:
    *   Arithmetic: `+`, `-`, `*`, `/`
    *   Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
    *   Logical: `&&`, `||` (with full short-circuiting)
    *   Postfix: `++`, `--`

### Standard Library Built-ins
*   **Formatting**: `print(...)`, `println(...)`
*   **Math**: `abs(x)`, `pow(base, exp)`, `sin(x)`, `cos(x)`, `tan(x)`
*   **Strings & Arrays**: `len(x)`, `contains(s, sub)`, `split(s, sep)`, `join(arr, sep)`, `to_upper(s)`, `to_lower(s)`
*   **Time**: `sleep(ms)`, `now()`
*   **Filesystem**: `read_file(path)`, `write_file(path, data)`

---

## Usage Instructions

### Compilation

Build the compiler/interpreter executable using Go:

```bash
go build -o rustx main.go
```

### Running Scripts

Execute scripts directly by passing the `.rx` file suffix:

```bash
./rustx samples/hello.rx
# Or via explicit run command:
./rustx run samples/factorial.rx
```

### Multiline REPL

Launch the interactive REPL shell:

```bash
./rustx repl
# Or simply:
./rustx
```

Example session:
```text
RustX Starter v0.1
>>> x := 5
>>> if x > 2 {
...     println("x is large");
... }
x is large
```

### Testing

Run the full, regression-free test suite:

```bash
go test ./...
```
