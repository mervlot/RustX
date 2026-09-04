# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** (`.rx`) is a minimal, clean, and scalable interpreted/compiled programming language foundation inspired by the elegance of Rust and the simplicity of Go.

Designed as a lightweight scripting and systems hybrid language, RustX Starter provides a modern foundation that is easy to write, run, and extend.

> **Note on Architecture:** Advanced memory management systems such as ownership, borrow checking, lifetimes, or generics are deliberately excluded from this foundation layer. The AST and evaluator are architected so that type checkers, borrow checkers, and bytecode/native compilers can be added in future stages without rewriting the core AST or interpreter.

---

## Language Features

- **Blended Rust & Go Syntax:** Modern syntax using `func` for function declarations, `var` and `:=` shorthand for variable bindings, curly brace blocks, and optional semicolons.
- **Data Types:** Primitive types (`int`, `float`, `bool`, `string`, `null`) and compound data structures (`array`).
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical short-circuiting (`&&`, `||`), and postfix operations (`++`, `--`).
- **Control Flow:** Flexible control structures including `if` / `else if` / `else`, `while`, and 3-part `for` loops with optional condition parenthesization.
- **Functions:** Higher-order function support, closures, optional parameter/return type annotations, and recursion.
- **Built-in Standard Library:** Standard utilities for printing (`print`, `println`), string manipulation (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem operations (`read_file`, `write_file`).
- **REPL & CLI:** Interactive multiline REPL, script runner (`rustx run file.rx`), syntax validator (`rustx build file.rx`), and direct file execution shorthand (`rustx file.rx`).

---

## Project Architecture

```
/
├── main.go               # CLI Entry point (run, repl, version, build, file shorthand)
├── token/                # Token types, definitions, and keyword mappings
├── lexer/                # Lexical scanner with line/col tracking and comment/escape support
├── ast/                  # Modular Abstract Syntax Tree nodes
├── parser/               # Recursive descent Pratt parser
├── object/               # Object runtime representation and environment scoping
├── evaluator/            # Core evaluator engine and standard library built-ins
├── errors/               # Human-readable diagnostic error formatters (SyntaxError, RuntimeError)
├── repl/                 # Multiline interactive REPL
└── samples/              # Sample RustX scripts (.rx)
```

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.18+ or higher.

### Building the Binary

Compile the `rustx` executable to your root directory:

```bash
go build -o rustx main.go
```

### Running Tests

Run all unit test suites:

```bash
go test ./...
```

---

## Usage & CLI

### 1. Interactive REPL

Start the RustX interactive prompt:

```bash
./rustx
# or
./rustx repl
```

Example session:

```rust
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(x, y))
30
```

### 2. Execute a Script

Run a `.rx` file:

```bash
./rustx run samples/factorial.rx
```

Or pass the file path directly:

```bash
./rustx samples/hello.rx
```

### 3. Syntax Verification

Validate syntax without executing runtime code:

```bash
./rustx build samples/control_flow.rx
```

### 4. Check Version

```bash
./rustx version
```

---

## Code Examples

### Example 1: Hello World

```rust
func main() {
    println("Hello World from RustX Starter!")
}
```

### Example 2: Recursive Factorial

```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    result := factorial(5)
    println("Factorial of 5 is:", result)
}
```

### Example 3: Variable Bindings & Loops

```rust
func main() {
    name := "RustX"
    if name == "RustX" {
        println("Language initialized:", name)
    }

    var count = 0
    while count < 3 {
        println("While loop count:", count)
        count++
    }

    for i := 0; i < 5; i++ {
        println("For loop i:", i)
    }
}
```

---

## Future Roadmap

The codebase is engineered to accommodate future language evolution:
1. Static type checking & type inference pass
2. Ownership rules, lifetimes, and borrow checker
3. Structs, Enums, and Pattern Matching
4. Modules and Import system
5. Async/Await and Concurrency primitives
6. Bytecode VM compiler and Native JIT/LLVM backend
