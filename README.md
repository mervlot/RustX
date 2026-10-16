<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

# RustX Starter

**RustX Starter** is a minimal, clean, scalable interpreted programming language designed as a lightweight scripting and systems hybrid. It blends syntax principles from **Rust** and **Go** to emphasize simplicity, speed, readability, and future extensibility.

> **Design Directive:** RustX Starter serves as the foundational architecture for future compiler evolution. Advanced memory safety mechanisms (borrow checker, ownership model, lifetimes) are deliberately deferred to ensure maximum flexibility while maintaining a scalable AST and runtime foundation.

---

## 🌟 Key Features

- **Blended Rust + Go Syntax:** `func` for functions, `var` for typed variables, `:=` for shorthand variable assignment, curly braces for scope, optional semicolons.
- **Dynamic & Static Hybrid Expressions:** Support for typed function signatures (e.g. `func add(a int, b int) int`), integers, floats, booleans, strings, and arrays.
- **Control Flow:** `if` / `else if` / `else`, `while`, and Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Operators:** Standard math (`+`, `-`, `*`, `/`), logical (`&&`, `||`, `!`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), and post-increment/decrement (`++`, `--`).
- **Standard Library Built-ins:**
  - **I/O:** `print`, `println`
  - **Strings:** `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math:** `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time:** `now`, `sleep`
  - **Filesystem:** `read_file`, `write_file`
- **Interactive REPL:** Multiline input support with bracket/brace balancing, error diagnostics, and persistent execution environment.
- **Extensible Architecture:** Recursive descent / Pratt parser emitting clean AST nodes structured for easy future extension to structs, traits, generics, modules, and native compilation / VM bytecode emission.

---

## 📁 Repository Structure

```text
rustx/
├── ast/            # AST node definitions (Statements, Expressions, Program)
├── errors/         # Standardized SyntaxError and RuntimeError diagnostic definitions
├── evaluator/      # Expression & statement evaluator, scoping, and stdlib built-ins
├── lexer/          # Lexical scanner / tokenizer supporting positions and comments
├── object/         # Runtime object definitions and environment bindings
├── parser/         # Recursive descent / Pratt precedence parser
├── repl/           # Interactive multiline REPL shell
├── samples/        # Example RustX scripts (.rx files)
├── token/          # Lexical token type constants and keyword mapping
├── main.go         # CLI entry point
├── go.mod          # Go module definition
└── README.md       # Project documentation
```

---

## 🚀 Getting Started & Building

### Prerequisites
- [Go](https://go.dev/) 1.20 or newer.

### Build the RustX Binary
From the project root:
```bash
go build -o rustx main.go
```

To run all unit tests:
```bash
go test ./...
```

---

## 💻 CLI Commands & Usage

### 1. Run a Script File
```bash
./rustx run samples/hello.rx
# Or shorthand direct execution:
./rustx samples/hello.rx
```

### 2. Interactive REPL Prompt
```bash
./rustx repl
# Or run without arguments:
./rustx
```

Example REPL Session:
```text
RustX Starter v0.1
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> func greet(name string) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### 3. Check Script Syntax
```bash
./rustx build samples/factorial.rx
# Output: Build successful: samples/factorial.rx checked with 0 errors.
```

### 4. Version
```bash
./rustx version
# Output: RustX Starter v0.1.0
```

---

## 📝 Syntax Examples

### Example 1: Factorial (Recursion)
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

### Example 2: Loops & Standard Library
```rustx
func main() {
    name := "RustX Starter"
    if name == "RustX Starter" {
        println("Language initialized successfully!")
    }

    sum := 0
    for i := 1; i <= 5; i++ {
        sum = sum + i
    }
    println("Sum of 1..5:", sum)

    message := "hello world"
    println("Upper:", to_upper(message))
}
```

---

## 🔬 Architecture Details & Future Roadmap

RustX Starter is designed in distinct, isolated layers:
1. **Lexer:** Stream-based tokenizer providing token location tracking (line, column) and comment stripping.
2. **Parser:** Pratt top-down operator precedence parser for expressions combined with recursive descent for statements. Semicolons are optional.
3. **AST:** Every AST node implements position tracking for precise error diagnostics.
4. **Interpreter / Evaluator:** Tree-walking evaluator operating on boxed runtime objects (`object.Object`) with lexically enclosed environments (`object.Environment`).

### Future Extension Architecture
The AST and Environment model are designed so future language extensions can be introduced incrementally:
- **Structs & Enums:** Can be added as AST statement nodes and evaluated as compound runtime objects.
- **Ownership & Borrow Checker:** Can be added as an intermediate AST analysis phase prior to evaluation or compilation.
- **Modules & Imports:** Can be handled by extending `Environment` to manage namespace resolution.
- **VM / Bytecode Compiler:** The existing AST can be fed into a bytecode compiler emitting instructions for a VM runtime without modifying front-end lexing and parsing.

---

## 📜 License
MIT License
