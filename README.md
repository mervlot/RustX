# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a minimal, clean, and extensible programming language inspired by **Rust** and **Go**. It blends the readability, concurrency-readiness, and systems elegance of Rust and Go into a lightweight scripting and systems foundation.

> **Design Notice:** In this initial starter phase, advanced memory safety systems (ownership, lifetimes, and borrow checking) are deliberately omitted. The architecture—from the recursive descent AST to the execution environment—is modularly designed so that a type checker, borrow checker, and native/bytecode compiler can be plugged in without refactoring the front-end parser and AST.

---

## Features

- **Blended Go & Rust Syntax**:
  - `func` for function declarations
  - `var` and `:=` for variable bindings
  - Optional semicolons
  - Curly-braced block scoping
  - Optional type annotations (`int`, `float`, `string`, `bool`) in function signatures and variable declarations
- **Rich Expressions & Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
  - Postfix increment & decrement: `++`, `--`
- **Control Flow**:
  - `if`, `else if`, `else`
  - `while` loops
  - `for` loops (both Go-style 3-part `for i := 0; i < 10; i++` and conditional `for condition`)
- **First-Class Functions & Recursion**:
  - Named and anonymous functions / closures
  - Recursion support
  - Automatic `main()` entrypoint invocation for script execution
- **Built-in Standard Library**:
  - **I/O**: `print`, `println`
  - **Strings**: `len`, `str_len`, `to_upper`, `to_lower`, `contains`, `split`
  - **Math**: `abs`, `sqrt`, `pow`, `min`, `max`
  - **Time**: `time_now`, `time_sleep`
  - **Filesystem**: `read_file`, `write_file`
  - **Data structures**: Dynamic arrays `[1, 2, 3]`, `append()`, and indexing `arr[0]`
- **Interactive REPL**:
  - Multiline statement entry with bracket balancing
  - Persistent state across statements
  - Human-friendly syntax and runtime error diagnostics
- **Clean CLI Tool**:
  - `rustx run <file.rx>`
  - `rustx <file.rx>` (shorthand execution)
  - `rustx repl`
  - `rustx build <file.rx>` (syntax & lexical verification)
  - `rustx version`

---

## Architecture

The project adheres to clean separation of concerns:

```
├── ast/         # Abstract Syntax Tree definitions (extensible for borrow checking & types)
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime evaluator and standard library built-ins
├── lexer/       # Lexical scanner with comment, string escape, and position tracking
├── object/      # Runtime value system and lexical environments
├── parser/      # Recursive descent and Pratt precedence parser
├── repl/        # Interactive REPL with multiline balance support
├── samples/     # Sample .rx programs demonstrating language features
├── token/       # Token specifications and keywords lookup
└── main.go      # CLI entrypoint
```

---

## Installation & Build

Ensure you have **Go 1.20+** installed:

```bash
# Clone the repository
git clone https://github.com/your-username/rustx.git
cd rustx

# Build the CLI binary
go build -o rustx main.go
```

To run all unit tests:
```bash
go test ./...
```

---

## Quick Start & Usage

### 1. Running .rx Script Files

```bash
./rustx run samples/01_hello.rx
# or shorthand:
./rustx samples/01_hello.rx
```

### 2. Interactive REPL

Run `./rustx` or `./rustx repl`:

```
RustX Starter v0.1
Type 'exit' or Ctrl+D to quit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func double(n int) int {
...     return n * 2
... }
>>> double(21)
42
```

### 3. Syntax Verification (Build)

Validate code without executing:
```bash
./rustx build samples/02_factorial.rx
```

---

## Code Examples

### Hello World (`main` entrypoint)
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Factorial (Recursion)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### Loops and Variables
```rustx
name := "RustX"
if name == "RustX" {
    println("Welcome to RustX")
}

for i := 0; i < 5; i++ {
    println("i =", i)
}
```

### Built-in Standard Library
```rustx
msg := "hello world"
println("Uppercase:", to_upper(msg))
println("Square root:", sqrt(64))

test_file := "/tmp/demo.txt"
write_file(test_file, "RustX I/O working!")
println(read_file(test_file))
```

---

## Future Roadmap

The codebase is engineered to scale directly into a fully fledged systems language:
1. **Static Type Checker**: Extend AST type annotations into static type analysis.
2. **Ownership & Borrow Checker**: Track lifetime bounds and borrow references (`&`, `&mut`) during semantic analysis.
3. **Structs & Traits**: Custom data structures, traits/interfaces, and methods.
4. **Bytecode VM / LLVM Native Compilation**: Replace AST tree-walk evaluation with bytecode or LLVM IR emission.
