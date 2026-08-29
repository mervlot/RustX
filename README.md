# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable foundation programming language inspired by **Rust** and **Go**. It blends the syntactic elegance of Rust (expressions, functions, control flow) with the simplicity and speed of Go.

> **Note on Architecture**: In this starter phase, RustX Starter behaves as a lightweight scripting + systems hybrid language. Advanced memory safety mechanisms like borrow checking, ownership, and lifetimes are intentionally omitted at this foundation level, but the AST and modular compiler architecture are explicitly designed to support them in future iterations without rewriting the engine.

---

## 🚀 Features

- **Blended Syntax**: `func` for function declarations, `var` and shorthand `:=` for variables, optional semicolons, optional condition parentheses.
- **Dynamic & Static Types**: Primitive support for `int`, `float`, `bool`, and `string`. Optional type annotations on parameters, variables, and return values.
- **Control Flow**: `if` / `else`, `while` loops, and Go-style 3-part `for` loops.
- **Standard Library**: Built-in functions for:
  - **Printing**: `print`, `println`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`
- **Interactive REPL**: Multiline block input tracking (`{}` and `()`), persistent environment state, and human-readable diagnostic messages.
- **CLI Commands**: Script execution, syntax checking (`build`), interactive prompt (`repl`), and version output.

---

## 📂 Architecture & Directory Structure

```text
rustx/
├── ast/        # Abstract Syntax Tree nodes for expressions and statements
├── errors/     # Custom SyntaxError and RuntimeError diagnostics with position tracking
├── evaluator/  # Interpreter, tree-walking AST evaluator, and standard library builtins
├── lexer/      # Lexical analyzer / Tokenizer supporting comments & escape sequences
├── object/     # Runtime Object model and Environment variable scoping
├── parser/     # Pratt / Recursive Descent parser
├── repl/       # Interactive REPL prompt handling multiline block inputs
├── samples/    # Sample .rx scripts showcasing language features
├── token/      # Token types, constants, and keyword lookup mappings
├── main.go     # CLI entry point
└── go.mod      # Go module specification
```

---

## 🛠️ Building & Running

### Prerequisites
- Go 1.20+ installed on your system.

### Build the Executable

```bash
go build -o rustx main.go
```

### Usage

1. **Run a `.rx` script**:
   ```bash
   ./rustx run samples/factorial.rx
   # Or directly:
   ./rustx samples/hello.rx
   ```

2. **Start Interactive REPL**:
   ```bash
   ./rustx repl
   # Or run without arguments:
   ./rustx
   ```
   ```text
   RustX Starter v0.1
   >>> x := 10
   >>> println(x * 2)
   20
   ```

3. **Check Syntax / Build**:
   ```bash
   ./rustx build samples/overview.rx
   ```

4. **Show Version**:
   ```bash
   ./rustx version
   ```

---

## 📝 Syntax Overview

### Variables & Assignment
```rustx
var count int = 10
name := "RustX"
count = count + 1
```

### Functions
```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
```

### Control Flow
```rustx
if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

for i := 0; i < 5; i++ {
    println(i)
}

var j = 5
while j > 0 {
    j--
}
```

---

## 🧪 Running Unit Tests

Run the full test suite across all packages:

```bash
go test ./...
```
