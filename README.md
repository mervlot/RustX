<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Starter Logo" width="200"/>
</p>

# RustX Starter

**RustX Starter** is a lightweight, compiled/interpreted scripting + systems hybrid programming language blending syntax and ergonomics inspired by Rust and Go.

It is designed with simplicity, speed, readability, and modular extensibility in mind as a foundation phase.

> **Note**: RustX Starter currently operates as a lightweight scripting engine. The architecture (AST, lexer, parser, runtime) is specifically structured so future language features (ownership, borrow checker, structs, traits, generics, modules, and native compilation) can be added cleanly without architectural refactoring.

---

## Key Features

- **Blended Rust & Go Syntax**:
  - `func` for function declarations
  - `var` for mutable variable declarations
  - `:=` shorthand assignment
  - Optional semicolons and optional block condition parentheses
- **Data Types**: `int`, `float`, `bool`, `string`, arrays (`[...]`)
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`
  - Logical: `&&`, `||`, `!`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Postfix: `++`, `--`
- **Control Flow**: `if`, `else`, `while`, `for` loops (3-part Go/C-style headers, single condition, and infinite loops)
- **Interactive REPL**: Multiline code block support, state preservation across evaluation prompts
- **Human-Readable Error Diagnostics**: Structured `SyntaxError` and `RuntimeError` reporting with line and column position tracking
- **Standard Library Built-ins**:
  - **Printing**: `print`, `println`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`

---

## Architecture Overview

```
/
├── ast/        # Abstract Syntax Tree nodes implementing Node, Statement, and Expression
├── errors/     # Structured SyntaxError and RuntimeError diagnostics with line/column
├── evaluator/  # AST Tree-walk interpreter and standard library built-ins
├── lexer/      # Lexical analyzer / tokenizer with comment & escape decoding
├── object/     # Runtime object system and lexical Environment scopes
├── parser/     # Pratt recursive descent parser
├── repl/       # Interactive REPL prompt with multiline buffer tracking
├── samples/    # Sample RustX (.rx) programs
├── token/      # Token types, constants, and keyword maps
├── main.go     # CLI entry point
└── go.mod      # Go module configuration
```

---

## Getting Started

### Building from Source

Ensure Go 1.20+ is installed on your system.

```bash
# Build the rustx executable
go build -o rustx main.go
```

### Running Scripts

Execute `.rx` scripts using `run` or shorthand command:

```bash
# Using the run command
./rustx run samples/hello.rx

# Shorthand invocation
./rustx samples/factorial.rx
```

### Interactive REPL

Start the REPL prompt:

```bash
./rustx repl
# or simply
./rustx
```

Example session:

```text
RustX Starter v0.1.0
Type "exit" or press Ctrl+C / Ctrl+D to exit.
>>> x := 5
>>> y := 10
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println("Sum:", add(x, y))
Sum: 15
>>> exit
```

### Syntactic Validation (`build`)

To validate lexical and syntactic correctness of a script without executing it:

```bash
./rustx build samples/loops.rx
```

---

## Examples

### Factorial Recursion (`samples/factorial.rx`)

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5:", factorial(5))
}
```

### Loops (`samples/loops.rx`)

```rustx
func main() {
    println("For loop counting:")
    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    println("While loop counting:")
    var count = 3
    while count > 0 {
        println("count =", count)
        count--
    }
}
```

---

## Testing

Run unit tests across all packages:

```bash
go test ./...
```
