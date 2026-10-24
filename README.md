# RustX Starter

**RustX Starter** is a minimal, clean, and scalable compiled/interpreted programming language inspired by **Rust** and **Go**. Designed with simplicity, speed, readability, and future extensibility in mind, RustX Starter behaves as a lightweight scripting + systems hybrid language.

---

## Architecture Overview

The project is structured with strict modularity to allow easy extensions (e.g. adding borrow checkers, lifetimes, generic types, modules, or a native bytecode VM compiler) without modifying core architectural abstractions.

```text
rustx/
├── token/       # Token types and keyword lookup definitions
├── errors/      # Human-readable SyntaxError and RuntimeError diagnostics
├── lexer/       # Tokenizer supporting strings, numbers, comments, and line/col tracking
├── ast/         # Abstract Syntax Tree nodes for statements and expressions
├── parser/      # Pratt recursive descent parser handling operators and control flow
├── object/      # Runtime evaluation environment and value objects
├── evaluator/   # AST interpreter with standard library built-ins
├── repl/        # Interactive REPL shell with multiline support
├── samples/     # Sample .rx demonstration programs
├── main.go      # CLI entry point (run, build, repl, version)
└── go.mod       # Go module declaration
```

---

## Key Features

- **Variables**: `var` declarations and `:=` shorthand assignments.
- **Data Types**: `int`, `float`, `bool`, `string`, `Array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Unary (`!`, `-`), Postfix (`++`, `--`).
- **Functions**: `func` with typed parameters, return types, recursive calling, and `main()` entry point auto-invocation.
- **Control Flow**: `if` / `else if` / `else`, `while`, and 3-part or single-condition `for` loops.
- **Comments**: Single-line `//` and multiline `/* ... */`.
- **Standard Library**:
  - Printing: `println()`, `print()`
  - Strings: `len()`, `to_upper()`, `to_lower()`, `split()`, `join()`, `contains()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - Time: `now()`, `sleep()`
  - Filesystem: `read_file()`, `write_file()`

---

## Installation & Building

### Prerequisites
- [Go 1.20+](https://golang.org/dl/)

### Build Executable
To compile the `rustx` executable binary:

```bash
go build -o rustx main.go
```

---

## CLI Usage

### 1. Execute a Script
```bash
./rustx run samples/hello.rx
```
Or shorthand:
```bash
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Launch the interactive prompt:
```bash
./rustx repl
```
Or simply:
```bash
./rustx
```

Example session:
```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func greet(name string) {
...     println("Hello", name)
... }
>>> greet("RustX")
Hello RustX
```

### 3. Check / Build Script
Perform lexical and syntactic validation without execution:
```bash
./rustx build samples/factorial.rx
```

---

## Code Examples

### Example 1: Hello World
```rustx
func main() {
    println("Hello World from RustX!")
}
```

### Example 2: Recursive Factorial
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### Example 3: Control Flow & Loops
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

var sum = 0
for i := 0; i < 5; i++ {
    sum = sum + i
}
println("Sum 0..4 =", sum)
```

---

## Running Tests

Run the full unit test suite across all packages:

```bash
go test ./...
```
