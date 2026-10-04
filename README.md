# RustX Starter

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. Designed for simplicity, execution speed, readability, and future extensibility, RustX Starter blends modern syntax aesthetics with lightweight scripting capabilities.

> **Note:** The current foundation phase intentionally omits complex ownership, lifetimes, and borrow checking systems. Its AST and evaluator architecture are designed so those systems can be seamlessly introduced in future iterations without rewriting the compiler core.

---

## 🚀 Key Features

- **Blended Rust & Go Syntax**: `func` for functions, `var` for mutable variables, `:=` shorthand declaration, optional semicolons, and curly brace scoping.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Control Flow**: Unparenthesized or parenthesized `if` / `else if` / `else`, `while`, and 3-part C-style `for` loops.
- **Functions & Recursion**: First-class functions, parameter type annotations, return statements, and full recursive support.
- **Standard Library Built-ins**:
  - **Output**: `print`, `println`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`
- **Interactive REPL**: Multiline block balancing, scope preservation across prompts, and human-readable diagnostics.
- **CLI Tooling**: Run scripts, validate file syntax, or launch interactive REPL sessions.

---

## 📂 Project Architecture

```
rustx/
├── ast/         # Abstract Syntax Tree nodes for statements & expressions
├── errors/      # Structured human-readable SyntaxError and RuntimeError
├── evaluator/   # Tree-walking evaluator and standard library built-ins
├── lexer/       # Lexical analyzer with position tracking and escape handling
├── object/      # Runtime object model and environment lexical scoping
├── parser/      # Pratt recursive descent parser
├── repl/        # Multiline interactive REPL prompt
├── samples/     # Sample .rx scripts showcasing features and syntax
├── token/       # Token definitions and keyword mappings
└── main.go      # CLI entry point
```

---

## 🛠️ Building & Installation

### Prerequisites
- Go 1.18+ installed on your system.

### Build Executable
```bash
go build -o rustx main.go
```

---

## 📖 Usage

### Running Scripts
Execute any `.rx` file directly:
```bash
./rustx run samples/hello.rx
```
Or use shorthand file execution:
```bash
./rustx samples/factorial.rx
```

### Interactive REPL
Start the interactive REPL:
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
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(x, y))
30
```

### Validate Syntax
Validate a script file's syntax without executing it:
```bash
./rustx build samples/features.rx
```

### Version Info
```bash
./rustx version
```

---

## 💻 Code Examples

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

func main() {
    println("Factorial of 5:", factorial(5))
}
```

### Example 3: Control Flow & Standard Library
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

for i := 0; i < 5; i++ {
    println("Iteration:", i)
}

text := "hello rustx"
println(to_upper(text))
```

---

## 🧪 Testing

Run all unit tests across lexer, parser, evaluator, and standard library components:
```bash
go test ./...
```
