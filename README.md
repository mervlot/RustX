# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a clean, scalable programming language inspired by **Rust** and **Go**, designed for simplicity, speed, readability, and future extensibility. It behaves as a lightweight hybrid scripting and systems foundation language with an extensible Pratt parser and evaluator architecture.

---

## Language Features

- **Variables**: `var` declarations, `:=` shorthand assignments, mutable by default.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`).
- **Functions**: Typed parameters and return values, recursion, first-class functions.
- **Control Flow**: `if` / `else if` / `else`, `while`, and 3-part `for` loops (parentheses optional).
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Built-in Standard Library**: `println`, `print`, `len`, `split`, `join`, `to_upper`, `to_lower`, `contains`, `abs`, `pow`, `sin`, `cos`, `tan`, `sleep`, `now`, `read_file`, `write_file`.

---

## Directory & Package Structure

```
.
├── ast/         # Abstract Syntax Tree nodes with token location tracking
├── errors/      # Human-readable SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime AST evaluator and Builtin standard library
├── lexer/       # Lexical analyzer with line & column tracking
├── object/      # Runtime Object system and scoped Environment handling
├── parser/      # Recursive descent / Pratt parser
├── repl/        # Interactive REPL with multiline support
├── samples/     # Example programs (.rx)
├── token/       # Token types and keyword mappings
├── main.go      # CLI entry point
└── README.md    # Documentation
```

---

## Architecture & Extensibility

RustX Starter is built with modularity in mind:

1. **Token & Position Tracking**: Every token records line and column information for precise error reporting.
2. **Pratt Parsing**: Expressions and precedence hierarchy are parsed using Pratt parsing rules, allowing clean grammar extension.
3. **AST Nodes**: All nodes implement `GetToken() token.Token` and contain parameter/return type fields to easily support future ownership checking, borrow checking, generics, or native code generation.
4. **Scoped Environment**: The runtime `Environment` supports lexical block scoping and variable resolution via recursive parent-scope lookups.

---

## Installation & CLI Usage

### Build Binary

```bash
go build -o rustx main.go
```

### Run Script File

```bash
./rustx run samples/hello.rx
# OR using shorthand execution
./rustx samples/hello.rx
```

### Interactive REPL

```bash
./rustx repl
# OR simply
./rustx
```

Example REPL session:
```
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+D to quit.
>>> x := 5
>>> println("Value of x is:", x)
Value of x is: 5
```

### Check Syntax (Build Command)

```bash
./rustx build samples/factorial.rx
```

### Print Version

```bash
./rustx version
```

---

## Unit Testing

Run the full test suite across all packages:

```bash
go test ./...
```
