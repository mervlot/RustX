# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by both **Rust** and **Go**. It blends the structural, expressive style of Rust with the structural simplicity and rapid-development model of Go, focusing on speed, readability, and clean architecture.

This version is designed as an extensible foundation, engineered so that advanced features (such as an ownership model, borrow checker, concurrency primitives, and async/await) can be introduced in the future without rewriting the compiler front-end or interpreter tree-walk components.

---

## Language Features

### 1. Simple, Expressive Syntax
- Standard variables with the `var` keyword and shorthand initialization with `:=`.
- Named and anonymous functions using the `func` keyword.
- Optional type annotations in function signatures (ignored during current phase evaluation but syntactically supported).
- Loop controls via unparenthesized `for` and `while` loop syntax.
- Standard logic control using `if` / `else`.
- Multi-line and single-line comment support.

### 2. Supported Data Types
- `int` (Integer, 64-bit precision)
- `float` (Float, 64-bit precision)
- `bool` (Boolean, `true` and `false`)
- `string` (Double-quoted with support for `\n`, `\t`, `\r`, `\"`, `\\` escape sequences)
- `array` (Zero-based indexed sequence with basic indexing assignment support, e.g., `arr[idx] = val`)

### 3. Comprehensive Operator Precedence
- Full arithmetic support (`+`, `-`, `*`, `/`)
- Relational operators (`==`, `!=`, `<`, `>`, `<=`, `>=`)
- Postfix operators for increment (`++`) and decrement (`--`) applied directly to identifiers.
- Short-circuiting logical operations with `&&` and `||`.

### 4. Rich Standard Library Built-ins
- **I/O & Logging**: `print`, `println`
- **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
- **String and Sequence manipulation**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
- **Time utilities**: `sleep` (milliseconds), `now` (Unix millisecond timestamp)
- **I/O and Filesystem**: `read_file` (reading paths), `write_file` (writing payloads)

---

## Directory Architecture

The compiler and interpreter runtime are structured cleanly in separate Go packages:

```
├── main.go             # Entrypoint & CLI execution
├── errors/             # Standardized Syntax and Runtime diagnostics
├── token/              # Lexical tokens and standard keywords
├── lexer/              # Finite automator / tokenizer
├── ast/                # Extensible Abstract Syntax Tree
├── parser/             # Recursive Descent & Pratt Pratt parser
├── object/             # Evaluation Value types and Scope Environments
├── evaluator/          # Evaluation runtime and standard library mappings
├── repl/               # multiline REPL execution loop
├── samples/            # Example .rx files showing off features
```

---

## Installation & Build Instructions

Prerequisites: **Go 1.24+**

To compile the binary executable from the root of the repository:

```bash
go build -o rustx main.go
```

To run all unit tests for lexing, parsing, and evaluation components:

```bash
go test ./...
```

---

## CLI Commands

### 1. Interactive REPL Prompt
Start an interactive interpreter session to run and debug expressions incrementally:
```bash
./rustx repl
# OR simply
./rustx
```

### 2. Execute a Script
Run an .rx script file directly:
```bash
./rustx run samples/hello.rx
# Or using the shorthand execution directly
./rustx samples/hello.rx
```

### 3. Syntax Verification (Dry Run Build)
Check the syntactic validity of a script without executing it:
```bash
./rustx build samples/factorial.rx
```

### 4. Check Language Tool Version
```bash
./rustx version
```

---

## Syntax Examples

### Variables & Shorthand Initialization
```rust
var language = "RustX"
version := 0.1
var isReady = true
```

### Unparenthesized Control Flow Loops & Shorthands
```rust
for i := 0; i < 5; i++ {
    println("Iteration number:", i)
}
```

### Recursive Function with Type Annotations
```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5)) // Outputs: 120
```
