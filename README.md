# RustX Starter

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntactic elegance and safety-first philosophy of Rust with the simplicity, speed, and readability of Go.

> **Note on Architecture & Future Readiness:**
> RustX Starter is designed as a foundational stepping-stone. It currently operates as a lightweight scripting + systems hybrid language. Advanced memory safety features like a **borrow checker, ownership model, and lifetimes are intentionally omitted in this phase**, but the AST and modular architecture are structured so these features (along with async/await, structs, enums, traits, modules, and native/bytecode compilation) can be seamlessly integrated in future iterations.

---

## Language Features

- **Blended Syntax:** Function syntax (`func`), mutable variable declarations (`var`), and shorthand initialization (`:=`).
- **Data Types:** `int` (64-bit integer), `float` (64-bit float), `bool`, `string`, `array`, `null`.
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`).
- **Control Flow:** Flexible `if` / `else` conditionals, `while` loops, and Go-style 3-part `for` loops (parentheses optional).
- **Functions:** Named and anonymous function support, type annotations, return values, and full recursion support. Auto-executes `main()` if defined in script mode.
- **Built-in Standard Library:**
  - **I/O & String Utilities:** `print`, `println`, `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`.
  - **Math Utilities:** `abs`, `pow`, `sin`, `cos`, `tan`.
  - **Time Utilities:** `now`, `sleep`.
  - **Filesystem Basics:** `read_file`, `write_file`.
  - **Array Utilities:** `push`, `pop`.
- **Comments:** Single-line (`//`) and multi-line (`/* ... */`).
- **Standardized Error Diagnostics:** Human-readable `SyntaxError` and `RuntimeError` messages with precise line and column location reporting.

---

## Code Example

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
    println("Factorial of 5 is:", factorial(5))
}
```

---

## Project Architecture

```
/ast         - Abstract Syntax Tree node definitions for statements and expressions
/errors      - Human-readable error diagnostics (SyntaxError, RuntimeError) with positions
/evaluator   - Tree-walking AST evaluator and built-in standard library functions
/lexer       - Lexical analyzer / tokenizer with position tracking and comment skipping
/object      - Object system representing runtime values and scoped environments
/parser      - Pratt parser (recursive descent) handling expressions, loops, and functions
/repl        - Multiline-capable interactive REPL prompt
/samples     - Sample .rx programs demonstrating language features
/token       - Token types, keyword lookup table, and location definitions
main.go      - CLI entry point supporting 'run', 'build', 'repl', 'version', and shorthand execution
```

---

## Build and Installation Instructions

### Prerequisites
- **Go 1.20+** installed on your system.

### Compiling the Binary
To compile the `rustx` executable:

```bash
go build -o rustx main.go
```

---

## CLI Usage

### 1. Interactive REPL
Launch the REPL by running `rustx` without arguments or using `rustx repl`:

```bash
./rustx repl
```

```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to exit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
```

### 2. Running Scripts
To run a `.rx` file:

```bash
./rustx run samples/hello.rx
```

Or using shorthand:

```bash
./rustx samples/factorial.rx
```

### 3. Syntactic Build Validation
Validate a `.rx` file's syntax without executing it:

```bash
./rustx build samples/loops.rx
```

### 4. Version Information
```bash
./rustx version
```

---

## Running Tests

Run the full Go unit test suite across all packages:

```bash
go test ./...
```
