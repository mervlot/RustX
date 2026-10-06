# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable scripting + systems hybrid programming language designed as an extensible foundation inspired by **Rust** and **Go**. It focuses on speed, simplicity, elegant syntax, and maintainability.

---

## Key Features

- **Blended Syntax**: Combines Rust and Go ergonomics (`func`, `var`, `:=`, optional semicolons, curly braces).
- **Core Types**: `int`, `float`, `bool`, `string`, `Array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), postfix (`++`, `--`), and shorthand assignment (`:=`).
- **Control Flow**: `if` / `else if` / `else`, `while`, and `for` loops (both single-condition and 3-part `for init; cond; post`).
- **First-Class Functions**: Top-level function declarations, closures, recursion, and optional type annotations on parameters and return types.
- **Human-Readable Diagnostics**: Clear syntax and runtime error messages with precise line and column numbers.
- **Standard Library Built-ins**:
  - **Strings & Collections**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `sleep`, `now`
  - **Filesystem**: `read_file`, `write_file`
- **Interactive REPL**: Interactive prompt supporting multiline input and state persistence across lines.

---

## Clean Architecture

The codebase follows a modular structure carefully designed for future extensions ( ownership, borrow checking, generics, structs, enums, bytecode VM, native compilation):

```
.
├── ast/         # AST interface and node definitions (position-aware)
├── errors/      # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walk evaluator engine and stdlib built-ins
├── lexer/       # Lexer and tokenizer supporting escape strings & comments
├── object/      # Object system, environment & lexical block scopes
├── parser/      # Pratt parser engine for precedence and unparenthesized control flow
├── repl/        # Interactive REPL prompt with multiline balance checking
├── samples/     # Sample RustX scripts (.rx)
├── token/       # Token types, keywords, and positions
├── main.go      # CLI entrypoint (run, repl, build, version)
└── README.md
```

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.18 or higher.

### Building RustX

Build the executable binary from the repository root:

```bash
go build -o rustx main.go
```

---

## CLI Usage

### Run a Script

```bash
./rustx run samples/hello.rx
```

Or pass a `.rx` file path directly:

```bash
./rustx samples/factorial.rx
```

### Start the REPL

```bash
./rustx repl
# or simply
./rustx
```

Example REPL session:

```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(10, 20))
30
```

### Check Syntax (Build Check)

To validate a script's syntax without executing it:

```bash
./rustx build samples/demo.rx
```

### Display Version

```bash
./rustx version
```

---

## Running Tests

Run the full unit test suite across all packages:

```bash
go test ./...
```

---

## Example Language Features

### Function Declaration & Recursion

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

### Variables & Loops

```rustx
var count = 0
name := "RustX Starter"

for i := 1; i <= 5; i++ {
    println("Step:", i)
}

k := 10
while k > 0 {
    k--
}
```

---

## Design & Future Compatibility

RustX Starter intentionally defers borrow checking, lifetimes, and advanced memory safety to maintain a lightweight foundation phase. However, every AST node retains token position and structure, making it straightforward to plug in a static analysis pass or borrow checker before code evaluation or generation in future releases.
