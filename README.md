<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Starter Logo" width="220" />
</p>

# RustX Starter

**RustX Starter** is a lightweight, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the ergonomics and simplicity of Go with the design philosophy, expressions, and future-readiness of Rust.

It serves as the foundation for a hybrid scripting and systems language, designed modularly so that advanced features (such as an ownership model, borrow checker, async/await, structs, and modules) can be integrated seamlessly in future versions.

---

## Features

- **Blended Go & Rust Syntax**:
  - `func` keyword for function declarations with optional type annotations
  - `var` declarations and `:=` shorthand assignments (mutable by default)
  - Control flow: `if` / `else if` / `else`, `while`, and `for` loops (Go-style 3-part headers or single conditions)
  - Pre/postfix increment and decrement (`++`, `--`)
  - Optional semicolons
  - Robust comments: single-line (`//`) and multiline (`/* ... */`)
- **Primitive & Composite Types**:
  - `int` (64-bit integers)
  - `float` (64-bit floating point)
  - `bool` (`true` and `false`)
  - `string` (with full escape sequence support: `\n`, `\t`, `\r`, `\"`, `\\`)
  - `array` (heterogeneous arrays, index expression & assignment support)
- **Functions & Recursion**:
  - First-class function support
  - Closures & lexical environment scoping
  - Full recursion support
- **Standard Library Built-ins**:
  - **I/O & Output**: `println(...)`, `print(...)`, `len(x)`
  - **Strings**: `str_to_upper(s)`, `str_to_lower(s)`, `str_contains(s, substr)`, `str_split(s, sep)`, `str_trim(s)`
  - **Math**: `math_sqrt(x)`, `math_pow(b, e)`, `math_abs(x)`, `math_floor(x)`, `math_ceil(x)`, `math_min(a, b)`, `math_max(a, b)`
  - **Time**: `time_now()`, `time_sleep(ms)`
  - **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Human-Readable Diagnostics**:
  - Formatted `SyntaxError` and `RuntimeError` with source line and column numbers. No panics on normal user errors.
- **Interactive Multiline REPL**:
  - State persistence across inputs, multiline block tracking with prompt (`>>>`) and continuation prompt (`...`).

---

## Architecture

The project is structured modularly:

```
├── ast/            # Abstract Syntax Tree node definitions
├── errors/         # Human-readable SyntaxError and RuntimeError diagnostics
├── evaluator/      # Tree-walking interpreter & standard library built-ins
├── lexer/          # Tokenizer tracking line/column positions and comments
├── object/         # Runtime object representation and lexical environment
├── parser/         # Recursive descent Pratt parser
├── repl/           # Interactive multiline REPL shell
├── samples/        # Example .rx programs
├── token/          # Lexical token definitions
└── main.go         # CLI entry point
```

### Future Extensibility
The AST and parser are designed with extensible node structures:
- Type annotations can be expanded into a full static type checker.
- Scoping in `object.Environment` and the AST can host an ownership tracking pass and borrow checker prior to evaluation.
- The tree-walking evaluator can be complemented or replaced with a bytecode compiler and VM without altering the front-end parser or AST.

---

## Installation & Build

### Prerequisites
- Go 1.20+ (Go 1.24+ recommended)

### Building from Source

```bash
# Clone the repository and navigate to root
cd rustx

# Build the executable
go build -o rustx main.go
```

---

## Usage

### 1. Execute a Script
```bash
./rustx run samples/hello.rx
# Or shorthand:
./rustx samples/hello.rx
```

### 2. Syntax Validation (Build Check)
```bash
./rustx build samples/factorial.rx
```

### 3. Interactive REPL
```bash
./rustx repl
# Or run with no arguments:
./rustx
```

Example REPL interaction:
```
RustX Starter v0.1
Type RustX expressions or statements. Press Ctrl+D or type exit to quit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func greet(name string) {
...     println("Hello", name)
... }
>>> greet("RustX")
Hello RustX
```

### 4. Version Information
```bash
./rustx version
```

---

## Code Examples

### Example 1: Hello World & Entry Point
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

println("Factorial of 5:", factorial(5))
```

### Example 3: Control Flow & Loops
```rustx
for i := 0; i < 5; i++ {
    println("i =", i)
}

count := 3
while count > 0 {
    println("countdown:", count)
    count--
}
```

---

## Running Tests

Run all unit tests across the entire codebase:

```bash
go test -v ./...
```
