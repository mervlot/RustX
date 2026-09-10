# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a lightweight, compiled/interpreted scripting + systems hybrid programming language inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and future extensibility, RustX Starter provides a clean foundation that can be expanded with advanced compiler systems (such as borrow checking, ownership models, generics, traits, and native compilation) without architectural overhauls.

---

## Key Features

- **Blended Rust & Go Syntax**: Clean block scoping with `{ ... }`, `func` for functions, `var` declarations, and `:=` shorthand assignment.
- **Dynamic + Primitive Types**: `int`, `float`, `bool`, and `string` primitives alongside `Array` objects.
- **Control Flow**: `if` / `else`, `while`, and 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library Built-ins**:
  - Printing: `println(...)`, `print(...)`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Filesystem: `read_file`, `write_file`
- **Interactive REPL**: Supports multiline block input, variable/function context persistence, and clean diagnostic error messages.
- **Robust Error Diagnostics**: Position-aware `SyntaxError` and `RuntimeError` reporting (`at line X column Y`).

---

## Architecture Overview

RustX Starter is designed with modular packages:

```
rustx/
├── token/          # Token type definitions, keywords, and line/col tracking
├── errors/         # Standardized SyntaxError and RuntimeError diagnostics
├── lexer/          # Lexical scanner for tokens, comments, and escape sequences
├── ast/            # AST node interfaces and concrete syntax tree representations
├── parser/         # Pratt recursive descent parser
├── object/         # Runtime type system, values, and lexical Environment
├── evaluator/      # AST interpreter engine and standard library built-ins
├── repl/           # Multiline REPL prompt and interactive execution loop
├── samples/        # Example .rx source files
└── main.go         # CLI interface (run, repl, version, build)
```

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/) version 1.22+ installed.

### Installation & Building

Clone the repository and build the binary:

```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

### 1. Run a Script
Execute a `.rx` script directly:

```bash
./rustx run samples/hello.rx
```

Or using shorthand file execution:

```bash
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Start the interactive shell:

```bash
./rustx repl
```

Or run `./rustx` without arguments:

```
RustX Starter v0.1 REPL
Type 'exit' or Ctrl+C to quit.
>>> x := 10
>>> func double(n int) int { return n * 2 }
>>> println(double(x))
20
```

### 3. Syntax Verification (Build)
Validate a `.rx` source file without executing runtime side-effects:

```bash
./rustx build samples/demo.rx
```

### 4. Version Check

```bash
./rustx version
```

---

## Syntax Examples

### Hello World

```rustx
func main() {
    println("Hello World")
}
```

### Recursion & Control Flow

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

### Loops & Variable Mutations

```rustx
var sum = 0
for i := 1; i <= 10; i++ {
    sum = sum + i
}
println("Sum 1..10 is:", sum)
```

---

## Testing

Run unit tests across all packages:

```bash
go test ./...
```

---

## Future Extensibility

RustX Starter is architected as an extensible foundation:
- **Ownership & Borrow Checker**: The AST nodes and Environment tracking can be enriched with lifetime and ownership metadata.
- **Generics, Structs & Enums**: Expressions and type descriptors are structured to allow user-defined data structures.
- **Native / Bytecode Compiler**: The AST representation decouples evaluation from parsing, enabling direct lowering to bytecode VM instructions or LLVM IR.
