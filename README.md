# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a lightweight, clean, and scalable scripting + systems hybrid programming language inspired by **Rust** and **Go**. It focuses on simplicity, fast execution, clean diagnostics, and future extensibility.

---

## 🚀 Key Features

- **Blended Syntax**: Integrates `func` functions, `var` declarations, `:=` shorthand assignments, unparenthesized control statements, and optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Control Flow**: `if` / `else if` / `else`, `while`, and Go-style 3-part `for` loops.
- **Built-in Standard Library**:
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `now`, `sleep`
  - Filesystem: `read_file`, `write_file`
- **Interactive REPL**: Multiline code block tracking, variable state preservation across evaluations, and human-readable position-aware diagnostics.
- **Extensible AST & Pratt Parser**: Designed with clear node interfaces and position tracking to easily incorporate future ownership models, borrow checkers, async/await, structs, traits, and module systems.

---

## 📂 Project Architecture

```text
/ast        - Extensible AST node definitions implementing GetToken()
/errors     - Human-readable SyntaxError and RuntimeError position-aware formatting
/evaluator  - Tree-walk interpreter, scoping rules, and standard library built-ins
/lexer      - Lexical analyzer supporting comments, escape sequences, and position tracking
/object     - Type-safe runtime object system and environment model
/parser     - Pratt / Recursive descent parser with operator precedence hierarchy
/repl       - Multiline interactive REPL shell
/samples    - Demonstrative RustX (.rx) scripts
/token      - Token types and keywords
main.go     - CLI entry point
go.mod      - Go module definition
```

---

## 🛠️ Installation & Building

Ensure Go 1.20+ is installed on your machine.

### Build Executable
```bash
go build -o rustx main.go
```

---

## 💻 Usage & CLI Commands

### 1. Run Script
Execute `.rx` scripts directly using either `run` or shorthand file path:
```bash
./rustx run samples/hello.rx
# OR
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Launch the interactive REPL:
```bash
./rustx repl
# OR
./rustx
```

Example REPL session:
```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> name := "RustX"
>>> println("Hello, " + name)
Hello, RustX
>>> for i := 0; i < 3; i++ {
...     println("i =", i)
... }
i = 0
i = 1
i = 2
```

### 3. Syntax Verification (Build Check)
Validate `.rx` script syntax without execution:
```bash
./rustx build samples/demo.rx
```

### 4. Check Version
```bash
./rustx version
```

---

## 📝 Syntax Overview & Code Examples

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    name := "RustX Starter"
    println("Welcome to", name)

    var sum = 0
    for i := 1; i <= 5; i++ {
        sum = sum + i
    }
    println("Sum of 1 to 5:", sum)
    println("Factorial of 5:", factorial(5))
}
```

---

## 🧪 Running Unit Tests

Run all unit tests across the codebase:
```bash
go test ./...
```
