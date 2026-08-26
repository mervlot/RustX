# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a minimal, clean, modular programming language designed as a blend between Rust and Go. It combines simple, modern syntax with high execution speed, clarity, and future extensibility.

Currently in its foundation phase, RustX Starter acts as a lightweight scripting and systems hybrid interpreter. It is architected so that advanced memory safety systems (ownership, borrow checking, lifetimes), concurrency models (`async`/`await`), structs, enums, traits, and native compilation can be added in future versions without rewriting the compiler or parser pipeline.

---

## Language Highlights

- **Rust + Go Hybrid Syntax**: `func` for functions, `var` for mutable variables, `:=` shorthand declaration, optional semicolons, optional condition parentheses.
- **Dynamic Expression Interpreter**: Fast Pratt parser, AST generation, and environment scoping.
- **Rich Standard Library**: Built-in functions for strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem I/O (`read_file`, `write_file`).
- **Interactive REPL**: Interactive prompt supporting multiline input balancing (`{` and `(`), state preservation, and clean error output.
- **Position-Aware Error Reporting**: Detailed human-readable `SyntaxError` and `RuntimeError` messages complete with line and column numbers.

---

## Directory Architecture

```
rustx/
├── main.go               # CLI entry point (run, build, repl, version)
├── go.mod                # Module definition (github.com/username/rustx)
├── token/                # Lexical tokens, literal constants, and line/column positions
├── errors/               # Structured SyntaxError & RuntimeError types
├── lexer/                # Tokenizer with comment filtering and escape decoding
├── ast/                  # Abstract Syntax Tree node interfaces and definitions
├── parser/               # Recursive descent Pratt parser implementation
├── object/               # Runtime object model & environment scoping
├── evaluator/            # Tree-walking evaluator & standard library built-ins
├── repl/                 # Multiline interactive REPL prompt
├── samples/              # Sample .rx demonstration scripts
└── assets/               # Branding assets
```

---

## Installation & Setup

Ensure you have **Go 1.20+** installed.

```bash
# Clone the repository
git clone https://github.com/username/rustx.git
cd rustx

# Build binary executable
go build -o rustx main.go
```

---

## Usage & Commands

### 1. Run Script Directly
```bash
./rustx run samples/hello.rx
# OR shorthand execution:
./rustx samples/hello.rx
```

### 2. Syntax Check (Build Mode)
```bash
./rustx build samples/factorial.rx
```

### 3. Interactive REPL
```bash
./rustx repl
```

**REPL Interaction Example:**
```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 5
>>> y := 10
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(x, y)
15
```

### 4. Version Check
```bash
./rustx version
```

---

## Code Examples

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
    println("Factorial of 5 is:", factorial(5))
}
```

### Example 3: Control Flow & Loops
```rustx
func main() {
    name := "RustX"
    if name == "RustX" {
        println("Language initialized:", name)
    }

    // While loop
    count := 0
    while count < 3 {
        println("Count:", count)
        count++
    }

    // For loop
    for i := 0; i < 5; i++ {
        println("Index:", i)
    }
}
```

---

## Architecture & Future Compatibility

RustX Starter is built around clean separation of concerns:
1. **Lexer (`lexer/`):** Scans source bytes into tokens while retaining source locations (`line`, `column`).
2. **Parser (`parser/`):** Utilizes Pratt parsing with precedence matching. AST nodes implement position interfaces for non-destructive future passes.
3. **AST (`ast/`):** Generic AST node interfaces (`Statement`, `Expression`) allow adding struct/enum definitions, type checking attributes, or borrow-checker annotations without breaking changes.
4. **Evaluator (`evaluator/`):** Uses an enclosed environment graph for variable scoping (`Get`, `Set`, `Update`).

**Future Extension Path:**
- Type checker pass prior to evaluation.
- Ownership & Borrow Checker static analyzer step between parsing and evaluation.
- Concurrency (`go` routines / `async`/`await`) primitives.
- Native LLVM / Cranelift compilation pipeline.

---

## Running Unit Tests

Run the full test suite across all packages:

```bash
go test ./...
```
