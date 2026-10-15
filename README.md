# RustX Starter

**RustX Starter** is a minimal, elegant, and extensible interpreted programming language inspired by **Rust** and **Go**. It blends the syntax and strengths of both languages to focus on simplicity, speed, readability, and future extensibility.

---

## Features

- **Blended Syntax:** Inspired by Go (`var`, `:=`, `func`, semicolons optional) and Rust (clean block structures and safety patterns).
- **Primitives & Types:** `int`, `float`, `bool`, `string`, arrays, dynamic runtime types.
- **Functions:** Named functions, anonymous closures, parameter type annotations, return type annotations, recursive support.
- **Control Flow:** `if` / `else if` / `else`, `while` loops, 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library:** Built-in support for:
  - **I/O & Printing:** `println()`, `print()`
  - **Strings:** `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math:** `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time:** `sleep()`, `now()`
  - **Filesystem:** `read_file()`, `write_file()`
- **Interactive REPL:** Multiline input buffering, persistent lexical scopes, line history.
- **Human-Readable Diagnostics:** Precise `SyntaxError` and `RuntimeError` line/column location tracking.
- **Extensible AST:** Designed so ownership, borrow checking, traits, modules, and bytecode compilation can be added in future stages without architecture rewrites.

---

## Directory Structure

```text
/rustx
├── ast/         # Abstract Syntax Tree definitions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime evaluator and built-in standard library
├── lexer/       # Lexer / Tokenizer with comment and escape handling
├── object/      # Runtime object system and environment scoping
├── parser/      # Pratt & Recursive Descent Parser
├── repl/        # Interactive REPL shell with multiline support
├── samples/     # Example RustX scripts (.rx)
├── token/       # Token types and definitions
├── go.mod       # Go module file
├── main.go      # CLI entry point
└── README.md    # Documentation and usage guide
```

---

## Installation & Building

Ensure you have **Go 1.18+** installed.

```bash
# Clone repository and build binary
go build -o rustx main.go
```

---

## Usage

### 1. Interactive REPL

Start the REPL by running `rustx` or `rustx repl`:

```bash
./rustx repl
```

Example session:
```text
RustX Starter REPL v0.1.0
Type 'exit' or press Ctrl+D to quit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
```

### 2. Run Scripts Directly

Execute any `.rx` script:

```bash
./rustx run samples/hello.rx
```

Or pass the script directly:

```bash
./rustx samples/factorial.rx
```

### 3. Syntax Validation (`build`)

Check syntax and lex errors in a file without executing runtime side effects:

```bash
./rustx build samples/loops_and_math.rx
```

### 4. Check Version

```bash
./rustx version
```

---

## Example Programs

### Basic Function & Recursion (`samples/factorial.rx`)

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    val := factorial(5)
    println("Factorial of 5 is:", val)
}
```

### Loops and Standard Library (`samples/loops_and_math.rx`)

```rustx
func main() {
    var sum = 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    println("Sum of 1..10 is:", sum)

    val := -42.5
    println("abs(-42.5):", abs(val))
    println("pow(2, 8):", pow(2, 8))
}
```

---

## Testing

Run unit test suites across all packages:

```bash
go test ./...
```

---

## Future Compatibility Roadmap

RustX Starter is built with clean layer isolation (Lexer -> Parser -> AST -> Evaluator). Future milestones include:
1. **Memory Model:** Ownership tracking, references (`&`), and compile-time borrow checker.
2. **Type System:** Structs, Enums, Generics, and Interface Traits.
3. **Concurrency:** Lightweight channels and async/await task execution.
4. **Target Compilation:** Bytecode Virtual Machine or Native LLVM backend compilation.
