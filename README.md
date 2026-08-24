# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180">
</p>

**RustX Starter** is a clean, minimal, modular programming language implementation inspired by **Rust** and **Go**. It combines elegant syntax (like `func`, `var`, `:=`, optional semicolons, and strong/simple typing) with high performance and human-readable diagnostic error reporting.

The architecture is deliberately modular to serve as a clean foundation for future extensions (such as ownership systems, borrow checkers, async/await, modules, structs, enums, traits, and native compilation or bytecode VMs).

---

## 🌟 Features

- **Blended Rust & Go Syntax**: `func` for functions, `var` and `:=` shorthand for variables, optional type annotations, and optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), logical (`&&`, `||`, `!`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), assignment (`=`, `:=`), and postfix increment/decrement (`++`, `--`).
- **Control Flow**: `if` / `else if` / `else`, `while` loops, and `for` loops (both 3-part `for init; cond; post` and single conditional `for cond`).
- **Standard Library Built-ins**:
  - **Printing**: `println()`, `print()`
  - **Strings**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time**: `sleep()`, `now()`
  - **Filesystem**: `read_file()`, `write_file()`
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Interactive REPL**: Interactive prompt with multiline buffer evaluation for unclosed `{` blocks or `(` parentheses.
- **Human-Readable Diagnostics**: Detailed line and column reporting for `SyntaxError` and `RuntimeError`.

---

## 🏗️ Architecture & Folder Structure

```
.
├── ast/         # Abstract Syntax Tree definitions (Node, Statement, Expression)
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── evaluator/   # AST Interpreter & Standard Library built-in functions
├── lexer/       # Tokenizer scanning source code into position-tracked tokens
├── object/      # Object system types, values, and Environment scoping table
├── parser/      # Recursive descent / Pratt expression parser
├── repl/        # Multiline REPL interactive shell
├── samples/     # Example .rx script files
├── token/       # Token constants and type lookup definitions
├── main.go      # CLI entry point for commands (run, repl, build, version)
└── README.md    # Documentation
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.20+ installed on your machine.

### Installation & Build

Build the `rustx` executable binary:

```bash
go build -o rustx main.go
```

Run test suites across all packages:

```bash
go test ./...
```

---

## 💻 Usage

### Run a Script

```bash
./rustx run samples/hello.rx
# OR shorthand execution
./rustx samples/factorial.rx
```

### Start the REPL

```bash
./rustx repl
# OR simply
./rustx
```

REPL Session Example:

```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to exit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> func fib(n int) int {
...     if n <= 1 { return n }
...     return fib(n - 1) + fib(n - 2)
... }
>>> println(fib(10))
55
```

### Check Script Syntax

```bash
./rustx build samples/demo.rx
```

---

## 📄 Example Programs

### 1. Hello World (`samples/hello.rx`)

```rustx
func main() {
    println("RustX Starter")
}
```

### 2. Recursive Factorial (`samples/factorial.rx`)

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### 3. Variables & Loops (`samples/demo.rx`)

```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

var age = 16
println("Age:", age)

println("Counting up to 5:")
for i := 1; i <= 5; i++ {
    println("i =", i)
}

msg := "rustx starter"
println("Uppercase:", to_upper(msg))
println("Length:", len(msg))
```

---

## 🔮 Future Roadmap

RustX Starter is designed so that future features can be integrated cleanly without rewriting the core AST or Parser:

- [ ] Ownership & Borrow Checker models
- [ ] Structs, Enums, and Pattern Matching
- [ ] Traits & Interfaces
- [ ] Concurrency (`go` / `async` / `await` paradigms)
- [ ] Bytecode VM & Native Compiler backend
- [ ] Module import system (`import "std/math"`)
