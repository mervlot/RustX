# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

---

## Features & Highlights

- **Blended Syntax**: Combines Go's clean control flow and declaration style (`func`, `var`, `:=`) with Rust-inspired syntax constructs.
- **Dynamic & Hybrid Scripting**: Flexible lightweight scripting feel without requiring full manual memory annotations in its foundation phase.
- **Extensible Architecture**: Modular Pratt parser and AST nodes designed for straightforward future additions (ownership, borrow checker, async, structs, traits, generics, native code / bytecode VM).
- **Rich Standard Library**: Built-in support for strings, math, time, and file system I/O.
- **Interactive REPL**: Interactive prompt with multiline block tracking (`{ ... }`, `( ... )`).
- **Clean Diagnostics**: Human-readable syntax errors and runtime diagnostics with line and column reporting.

---

## Directory Structure

```text
rustx/
├── ast/         # Abstract Syntax Tree node definitions
├── errors/      # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/   # Interpreter, runtime evaluation engine, and stdlib built-ins
├── lexer/       # Lexical analyzer / tokenizer with position tracking
├── object/      # Object system, type representations, and scoped environment
├── parser/      # Recursive descent Pratt parser with operator precedence
├── repl/        # Interactive REPL prompt with multiline support
├── samples/     # Sample programs demonstrating language features
├── token/       # Token specifications and keywords
├── main.go      # CLI entry point
├── go.mod       # Go module file
└── README.md    # Documentation
```

---

## Installation & CLI Usage

### Build from Source

Ensure Go 1.20+ is installed:

```bash
go build -o rustx main.go
```

### Commands

#### Interactive REPL

```bash
./rustx repl
# Or simply
./rustx
```

```text
RustX Starter v0.1.0 Interactive REPL
Type 'exit' or 'quit' to exit.
>>> x := 10
>>> func double(n int) int { return n * 2; }
>>> double(x)
20
```

#### Run a Script

```bash
./rustx run samples/hello.rx
# Or shorthand
./rustx samples/hello.rx
```

#### Build / Validate Syntax

```bash
./rustx build samples/hello.rx
```

#### Version Info

```bash
./rustx version
```

---

## Language Syntax Guide

### Variables & Shorthand

```rustx
var name = "RustX"
var count int = 100
age := 16
```

### Data Types

- `int`: 64-bit signed integers (`42`)
- `float`: IEEE 754 64-bit floats (`3.14159`)
- `bool`: Booleans (`true`, `false`)
- `string`: String literals with escape sequences (`"Hello\nWorld"`)
- `array`: Lists (`[1, 2, 3, 4]`)

### Control Flow

```rustx
if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

for i := 0; i < 5; i++ {
    println("Count:", i)
}

while count > 0 {
    count--
}
```

### Functions

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

---

## Standard Library Built-ins

| Function | Description |
|---|---|
| `println(...)` | Prints space-separated values followed by a newline |
| `print(...)` | Prints space-separated values without a newline |
| `len(val)` | Returns the length of a string or array |
| `contains(s, sub)` | Returns true if string contains substring |
| `split(s, sep)` | Splits string by separator into an array |
| `join(arr, sep)` | Joins array of elements into a string using separator |
| `to_upper(s)` | Converts string to uppercase |
| `to_lower(s)` | Converts string to lowercase |
| `abs(n)` | Returns the absolute value of integer or float |
| `pow(base, exp)` | Calculates base raised to exponent |
| `sin(x)` / `cos(x)` / `tan(x)` | Trigonometric functions |
| `sleep(ms)` | Pauses execution for given milliseconds |
| `now()` | Returns Unix epoch timestamp in milliseconds |
| `read_file(path)` | Reads file contents into a string |
| `write_file(path, content)` | Writes string content to specified filepath |

---

## Testing

Run all unit tests across lexer, parser, and evaluator:

```bash
go test ./...
```

---

## Future Roadmap & Extensibility

RustX Starter was explicitly designed so advanced systems can be layered onto this foundation without major parser rewrites:

1. **Ownership & Borrow Checker**: Linear type analysis pass over AST nodes before evaluation or compilation.
2. **Structs, Enums & Traits**: Extending Pratt parser production rules and Object system.
3. **Concurrency & Async**: Adding `async`/`await` primitives and lightweight channels.
4. **Bytecode VM & Native Compiler**: Target intermediate representation (IR) or LLVM / Cranelift backend.
