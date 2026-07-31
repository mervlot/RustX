# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. This repository provides a complete foundation-level implementation of the language written in Go.

## Architectural Overview

The project is built on a highly modular architecture that decouples each stage of processing to make the language easy to extend:

- **token**: Defines the set of lexemes, keywords, and structural elements of the language, maintaining position tracking (lines and columns) for better errors.
- **errors**: Provides structured syntax (`SyntaxError`) and runtime (`RuntimeError`) diagnostics with position-aware formatting.
- **lexer**: Converts raw source text into a stream of tokens, decoding backslash escapes, skipping multiline comment blocks, and handling integers/floats.
- **ast**: Structures the syntax trees implementing the position-aware Node interface.
- **parser**: A Pratt parser combined with Recursive Descent, supporting unparenthesized conditions and custom Go-style loops.
- **object**: The evaluation runtime object model supporting integers, floats, booleans, strings, arrays, custom functions, environment frame variables, and recursive variable updates.
- **evaluator**: A tree-walking interpreter mapping AST nodes to objects, with short-circuiting logic, error checking (e.g. division by zero), and scoping rules.
- **repl**: An interactive REPL prompt that handles balanced brace checking to seamlessly support multiline declarations.

---

## Installation & Building

Make sure you have Go installed (Go 1.18+ recommended).

1. Build the compiler binary from the repository root:
   ```bash
   go build -o rustx main.go
   ```

2. Run the tests to ensure everything is correct:
   ```bash
   go test ./...
   ```

---

## Usage & CLI Commands

RustX Starter provides a complete command-line interface.

### Running a script directly

```bash
./rustx run samples/factorial.rx
```

You can also use the shorthand direct execution syntax:
```bash
./rustx samples/factorial.rx
```

### Interactive REPL Prompt

Launch the interactive prompt:
```bash
./rustx repl
```
or simply:
```bash
./rustx
```

Within the REPL:
```
RustX Starter v0.1
Type commands or 'exit' to quit.
>>> x := 5
>>> println(x)
5
>>> func add(a, b) {
...     return a + b;
... }
>>> add(x, 10)
15
```

### Displaying Version

```bash
./rustx version
```

---

## Language Features

### 1. Variables
Declared using the `var` keyword or shorthand `:=` declaration:
```rust
var x = 10;
y := 20;
```

### 2. Standard Types
Integers, Floats, Booleans, Strings, and Arrays:
```rust
arr := [1, 2.5, "hello", true]
```

### 3. Control Flow
Supports standard `if-else` blocks, `while` loops, and Go-style `for` loops (parentheses are optional):
```rust
for i := 0; i < 5; i++ {
    println(i)
}
```

### 4. Functions & Recursion
Anonymous functions or Named Functions are supported:
```rust
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
```

### 5. Built-in Functions (Standard Library)
- **Printing**: `print`, `println`
- **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
- **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
- **Time**: `sleep`, `now`
- **Filesystem**: `read_file`, `write_file`
