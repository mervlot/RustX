# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a lightweight, compiled/interpreted scripting and systems hybrid programming language designed to seamlessly blend the clarity, performance, and modern ergonomics of **Rust** and **Go**.

---

## Language Overview & Syntax

RustX Starter features an expressive syntax with simple typing, optional semicolons, and clean block scoping.

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    name := "RustX Starter"
    var age = 16

    println("Welcome to", name)
    if age >= 16 {
        println("Factorial of 5 is:", factorial(5))
    }
}
```

### Core Features
- **Variables**: `var x = 10` or shorthand declaration `x := 10`.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`).
- **Control Flow**: `if` / `else if` / `else`, `while`, and Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Functions**: Named functions, anonymous functions/closures, recursion, typed parameters, and return types.
- **Built-in Functions**: `print`, `println`, `len`, `str_contains`, `str_upper`, `str_lower`, `math_sqrt`, `math_pow`, `math_abs`, `time_now`, `read_file`, `write_file`.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).

---

## Directory & Architecture Structure

The codebase is built with modularity and extensibility in mind to support future compiler passes (e.g., borrow checker, ownership analyzer, bytecode VM, LLVM backend):

```
.
├── ast/         # Abstract Syntax Tree nodes with position metadata
├── errors/      # Structured diagnostic errors with line and column reporting
├── evaluator/   # Tree-walk interpreter and standard library built-ins
├── lexer/       # Tokenizer with line/column tracking & comment processing
├── object/      # Runtime object system and lexical environment
├── parser/      # Pratt recursive-descent parser
├── repl/        # Interactive REPL prompt with multiline block tracking
├── samples/     # Example programs (.rx)
├── token/       # Token type definitions and keyword dictionary
├── main.go      # CLI entry point
└── go.mod       # Go module declaration
```

---

## CLI Usage & Commands

Build the `rustx` executable:

```bash
go build -o rustx main.go
```

### Run Scripts
Execute a RustX script:
```bash
./rustx run samples/hello.rx
# OR shorthand execution:
./rustx samples/factorial.rx
```

### Interactive REPL Prompt
Start the interactive REPL:
```bash
./rustx repl
# OR simply:
./rustx
```

Example session:
```
RustX Starter v0.1 Interactive REPL
Type 'exit' or 'quit' to exit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> func greet(name) {
...     return "Hello " + name
... }
>>> greet("World")
Hello World
```

### Validate Syntax (Build Check)
Check a script for syntax/lexical errors without execution:
```bash
./rustx build samples/factorial.rx
```

### Display Version
```bash
./rustx version
```

---

## Testing

Run all unit tests across the repository:

```bash
go test ./...
```
