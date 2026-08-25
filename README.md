# RustX Starter

**RustX Starter** is a minimal, elegant, modular compiled/interpreted programming language inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and future extensibility, it blends modern syntax with robust execution.

---

## Language Highlights

- **Hybrid Syntax**: Combines `func`, `var`, shorthand `:=`, typed parameter definitions, optional semicolons, and curly-brace control blocks.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`, `null`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`).
- **Control Flow**: `if`, `else if`, `else`, `while`, `for` (Go-style 3-part loops).
- **Functions**: Parameter type annotations, return type signatures, named functions, recursion, auto-execution of top-level `main()`.
- **Standard Library**:
  - Printing: `print()`, `println()`
  - Strings & Arrays: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - System: `sleep()`, `now()`, `read_file()`, `write_file()`
- **Extensible Architecture**: Structured with line/column position reporting and clean AST separation so memory safety models (borrow checking, ownership), concurrency, generics, and native/bytecode compilation can be integrated seamlessly in future versions.

---

## Directory Structure

```text
.
├── ast/         # Abstract Syntax Tree interfaces and node definitions
├── errors/      # Custom SyntaxError and RuntimeError diagnostics with line/column positioning
├── evaluator/   # Runtime interpreter and built-in standard library implementation
├── lexer/       # Tokenizer with comment and escape sequence support
├── object/      # Runtime Object model and Environment scope management
├── parser/      # Pratt recursive-descent parser
├── repl/        # Interactive multiline REPL prompt
├── samples/     # Sample RustX (.rx) programs
├── token/       # Token definitions and keyword lookup table
├── main.go      # CLI entry point
├── go.mod       # Go module configuration
└── README.md    # Documentation and usage guide
```

---

## Getting Started

### Prerequisites
- [Go 1.20+](https://go.dev/dl/)

### Build Binary

To compile the `rustx` executable:

```bash
go build -o rustx main.go
```

---

## Command Line Usage

### 1. Execute Script Files
Run `.rx` files directly using either `rustx run` or shorthand path execution:

```bash
./rustx run samples/hello.rx
# OR
./rustx samples/hello.rx
```

### 2. Interactive REPL Prompt
Launch the multiline interactive prompt:

```bash
./rustx repl
# OR simply
./rustx
```

Example REPL interaction:
```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> exit
```

### 3. Check Version
```bash
./rustx version
```

### 4. Build / Syntax Validation
Check scripts for syntax errors without executing:

```bash
./rustx build samples/factorial.rx
```

---

## Example Programs

### Example 1: Hello World
```rustx
func main() {
    println("Hello World from RustX Starter!")
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

### Example 3: Control Flow & Standard Library
```rustx
name := "RustX Starter"

if name == "RustX Starter" {
    println("Language initialized successfully.")
}

var sum = 0
for i := 0; i < 5; i++ {
    sum = sum + i
}

println("Sum of numbers from 0 to 4 is:", sum)

str := "rustx programming language"
println("Uppercase:", to_upper(str))
```

---

## Testing

Run all unit tests across lexer, parser, and evaluator modules:

```bash
go test ./...
```
