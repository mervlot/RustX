# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a clean, modern, lightweight compiled/interpreted programming language designed by blending the elegance and power of **Rust** and **Go**.

It focuses on **simplicity**, **speed**, **readability**, and **future extensibility**.

---

## 🚀 Key Features

- **Blended Rust + Go Syntax**: Clean block structure using `func`, `var`, shorthand `:=`, optional type annotations, and optional semicolons.
- **Dynamic & Direct Execution**: Run `.rx` scripts directly or interactively via REPL.
- **Control Flow**: Unparenthesized `if`, `else`, `while`, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library Built-ins**:
  - **Printing**: `print`, `println`
  - **String utilities**: `len`, `split`, `join`, `to_upper`, `to_lower`, `contains`
  - **Math utilities**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time & Filesystem**: `now`, `sleep`, `read_file`, `write_file`
- **Future-Ready Architecture**: Pratt parser and AST design structured so advanced systems (ownership model, borrow checker, generics, modules, async) can be added without rewriting the core runtime.

---

## 🛠 Directory Architecture

```text
rustx/
├── token/       # Token definitions and keyword lookup
├── lexer/       # Lexical scanner supporting string escapes & comments
├── ast/         # AST node implementations with source position tracking
├── errors/      # Human-readable SyntaxError and RuntimeError diagnostics
├── parser/      # Pratt recursive descent parser
├── object/      # Runtime object system and scoping Environment
├── evaluator/   # Evaluator engine and standard library built-ins
├── repl/        # Interactive REPL with multi-line block detection
├── samples/     # Sample RustX (.rx) programs
├── main.go      # CLI entry point
└── README.md    # Documentation
```

---

## 💻 Building and Running

### Prerequisites
- **Go** 1.18+ or later

### Build Binary
```bash
go build -o rustx main.go
```

### Run Tests
```bash
go test ./...
```

---

## 📖 Usage & CLI

```bash
# Start Interactive REPL
./rustx
# or
./rustx repl

# Run a RustX Script
./rustx run samples/01_hello.rx

# Shorthand Script Execution
./rustx samples/02_factorial.rx

# Syntax Verification (Build)
./rustx build samples/03_loops.rx

# Show Version
./rustx version
```

---

## 📝 Syntax Overview

### 1. Variables & Types
```rustx
var name = "RustX"
age := 20
var active bool = true
```

### 2. Functions & Recursion
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

### 3. Control Flow
```rustx
// For Loop
for i := 0; i < 5; i++ {
    println("i =", i)
}

// While Loop
count := 3
while count > 0 {
    println("Countdown:", count)
    count--
}
```

---

## 🔮 Future Extensibility

RustX Starter is designed as a foundational stepping stone. Its AST, Environment object model, and Pratt parser engine are built to cleanly accommodate upcoming additions:
- Borrow checker & Ownership model
- Structs, Enums, Traits/Interfaces
- Module & Import system
- Concurrency & Async/Await
- Native Code compilation / Bytecode VM target

---

## 📄 License

MIT License.
