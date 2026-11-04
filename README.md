# RustX Starter

![RustX](assets/rustx-logo.png)

A clean, minimal, and scalable interpreted and compiled programming language foundation inspired naturally by **Rust** and **Go**. RustX Starter blends ergonomic syntax, lightweight semantics, and high readability while establishing an extensible architecture prepared for future systems features (e.g. ownership, borrow checker, structs, enums, concurrency).

---

## Architecture Overview

RustX Starter follows a modular, decoupled compiler & interpreter pipeline:

```
Source (.rx)
    │
    ▼
[Lexer] ──────────► Stream of Tokens (with Line/Col tracking)
    │
    ▼
[Parser] ─────────► Abstract Syntax Tree (Pratt Recursive Descent)
    │
    ▼
[Evaluator] ──────► Object Runtime Environment & Standard Library
```

### Directory Structure

- `token/`: Token types, keywords, operators, delimiters, and source position definitions.
- `lexer/`: Lexical scanner supporting single-line comments (`//`), multiline comments (`/* */`), string escapes, numbers, and identifiers.
- `ast/`: AST node hierarchies for expressions, statements, function literals, type annotations, and control structures.
- `parser/`: Pratt top-down operator precedence / recursive descent parser with clean syntax error reporting.
- `object/`: Runtime value representations (`Integer`, `Float`, `Boolean`, `String`, `Array`, `Function`, `Builtin`, `Error`, `Null`) and scoped environments.
- `evaluator/`: Tree-walking runtime evaluation engine and built-in standard library functions (`builtins.go`).
- `repl/`: Interactive multiline REPL with persistent environments and balanced delimiter tracking.
- `errors/`: Structured `SyntaxError` and `RuntimeError` types with location context.
- `samples/`: Sample RustX programs illustrating language constructs and features.
- `main.go`: Command-line interface entry point.

---

## Language Syntax & Features

### 1. Variables & Assignment
```rustx
var count = 10
var message string = "Welcome to RustX"
shorthand := 42
count = count + 1
```

### 2. Data Types
- `int`: 64-bit signed integers (`42`, `-7`)
- `float`: 64-bit IEEE floating-point numbers (`3.14`, `-0.5`)
- `bool`: `true` or `false`
- `string`: double-quoted UTF-8 strings (`"Hello\nWorld"`)
- `array`: dynamic collections (`[1, 2, 3]`)

### 3. Operators
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Postfix: `++`, `--`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `&&`, `||`, `!`

### 4. Functions & Recursion
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial:", factorial(5))
}
```

### 5. Control Flow
```rustx
// If - Else
if x > 10 {
    println("Greater")
} else {
    println("Smaller or equal")
}

// For Loop (Go-style 3-part or condition)
for i := 0; i < 5; i++ {
    println(i)
}

// While Loop
count := 3
while count > 0 {
    println(count)
    count--
}
```

### 6. Standard Library Built-ins
- **I/O & Print**: `println(...)`, `print(...)`
- **Strings**: `len(s)`, `to_upper(s)`, `to_lower(s)`, `str_contains(s, sub)`
- **Math**: `sqrt(x)`, `abs(x)`, `pow(b, e)`
- **Time**: `time_now()`, `sleep_ms(ms)`
- **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Collections**: `push(array, elem)`

---

## CLI & Usage

### Building RustX Binary
```bash
go build -o rustx main.go
```

### Running Scripts
```bash
# Execute with run command
./rustx run samples/hello.rx

# Or shorthand
./rustx samples/factorial.rx
```

### Interactive REPL
```bash
./rustx repl
# or simply
./rustx
```

REPL session:
```
RustX Starter v0.1
Type 'exit' or press Ctrl+C/Ctrl+D to exit.
>>> x := 5
>>> y := 10
>>> println("Result:", x + y)
Result: 15
```

### Syntax Validation / Build Check
```bash
./rustx build samples/loops.rx
```

### Running Tests
```bash
go test -v ./...
```
