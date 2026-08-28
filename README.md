<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

# RustX Starter

RustX Starter is a clean, modular, systems-inspired programming language interpreter blending syntax from **Rust** and **Go**. Designed as a light scripting + systems hybrid foundation, it prioritizes speed, readability, and clean architecture so advanced features (like ownership, borrow checking, and static compilation) can be added in future stages without architectural rewrites.

---

## Language Features

- **Blended Rust + Go Syntax**: `func` for functions, `var` and `:=` for variable declarations, postfix `++`/`--`, optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), and short-circuit logical operators (`&&`, `||`).
- **Control Flow**: `if` / `else if` / `else`, `while`, and Go-style `for` loops.
- **Standard Library**:
  - Printing: `print()`, `println()`
  - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - Time: `sleep()`, `now()`
  - Filesystem: `read_file()`, `write_file()`
- **Interactive REPL**: Multiline code buffer tracking, state persistence, clean diagnostic error output.

---

## Directory Architecture

```
rustx/
├── ast/            # Abstract Syntax Tree nodes for statements & expressions
├── errors/         # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/      # Runtime evaluation logic & standard library built-in bindings
├── lexer/          # Tokenizer with line & column position tracking
├── object/         # Object system and Environment scoping model
├── parser/         # Pratt recursive descent parser
├── repl/           # Multiline interactive REPL prompt
├── samples/        # Example RustX (.rx) programs
├── token/          # Token definitions and keyword lookup
├── main.go         # CLI entry point
└── go.mod          # Go module file
```

---

## Installation & Build Instructions

Prerequisites: Go 1.20+

### Building the Executable

```bash
# Build the rustx binary in the root directory
go build -o rustx main.go
```

---

## Usage

### 1. Interactive REPL
```bash
./rustx repl
# OR simply
./rustx
```

Example REPL session:
```
RustX Starter v0.1.0 Interactive REPL
Type 'exit' or press Ctrl+D to exit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(10, 20)
30
```

### 2. Running a Script
```bash
./rustx run samples/hello.rx
# OR shorthand
./rustx samples/factorial.rx
```

### 3. Syntax Validation / Build Check
```bash
./rustx build samples/control_flow.rx
```

### 4. Version Check
```bash
./rustx version
```

---

## Code Examples

### Factorial (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5 is:", factorial(5))
}
```

### Control Flow (`samples/control_flow.rx`)
```rustx
func main() {
    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    var count = 3
    while count > 0 {
        println("count =", count)
        count--
    }
}
```

---

## Testing

Run all unit tests across packages:
```bash
go test ./...
```
