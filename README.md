<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

# RustX Starter

**RustX Starter** is a clean, scalable programming language blending the elegance of **Go** and **Rust**. It combines simplicity, readable syntax, and speed into a lightweight scripting and systems-inspired language designed for future extensibility.

---

## Highlights

- **Syntax Blend**: Intuitive syntax using `func`, `var`, `:=`, curly braces `{ ... }`, and optional semicolons.
- **Modern Control Flow**: Support for `if`/`else`, `while`, and 3-part Go-style `for` loops (`for i := 0; i < 10; i++ { ... }`).
- **Interactive REPL**: Interactive prompt supporting multiline input, statement completion, and persistent session state.
- **Diagnostic Errors**: Human-readable syntax and runtime error messages with precise line and column numbers.
- **Built-in Standard Library**: First-class support for strings, math, time, arrays, and filesystem I/O.
- **Designed for Growth**: AST and runtime architecture decoupled and ready for future borrow checking, ownership models, generics, structs, and VM execution.

---

## Project Architecture

```
rustx/
├── token/        # Token definitions and keyword lookup
├── errors/       # Standardized SyntaxError and RuntimeError diagnostics
├── lexer/        # Lexical analyzer with whitespace/comment handling and position tracking
├── ast/          # Abstract Syntax Tree nodes implementing Node, Statement, and Expression
├── parser/       # Pratt recursive descent parser
├── object/       # Runtime values, type definitions, and lexical Environment
├── evaluator/    # AST tree-walking interpreter and built-in standard library
├── repl/         # Interactive REPL with multiline buffering
├── samples/      # Example .rx programs demonstrating language features
├── main.go       # CLI entry point (run, build, repl, version)
└── README.md     # Documentation and specifications
```

---

## Language Syntax

### 1. Variables and Types

Variables can be declared using `var` (with optional type annotations) or shorthand `:=`:

```rustx
var x = 10
var message string = "Hello RustX"
name := "Muhammad"
isReady := true
pi := 3.14159
```

Supported primitive data types:
- `int` (64-bit integer)
- `float` (64-bit floating point)
- `bool` (`true`, `false`)
- `string` (double-quoted with escape characters `\n`, `\t`, `\"`, `\\`)
- `array` (`[1, 2, "three"]`)

### 2. Operators

- **Arithmetic**: `+`, `-`, `*`, `/`, `%`
- **Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Logical**: `&&`, `||`, `!`
- **Postfix**: `++`, `--`

### 3. Control Flow

#### `if` / `else`
Parentheses are optional:

```rustx
if score >= 90 {
    println("Grade: A")
} else if score >= 80 {
    println("Grade: B")
} else {
    println("Grade: C")
}
```

#### `while` loops
```rustx
var n = 3
while n > 0 {
    println(n)
    n--
}
```

#### `for` loops
Full 3-part C/Go-style loop:

```rustx
for i := 0; i < 5; i++ {
    println(i)
}
```

### 4. Functions & Recursion

Functions are declared with `func`. Parameter and return type annotations are optional or strong:

```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5)) // prints 120
```

If a zero-argument `func main() { ... }` is present in a script file, RustX will invoke it as the entry point automatically.

---

## Standard Library Built-ins

| Category | Functions |
|---|---|
| **I/O** | `print(...)`, `println(...)` |
| **Strings** | `len(s)`, `to_upper(s)`, `to_lower(s)`, `contains(s, sub)`, `split(s, sep)`, `join(arr, sep)`, `trim(s)` |
| **Math** | `abs(x)`, `sqrt(x)`, `min(a, b)`, `max(a, b)`, `pow(x, y)`, `floor(x)`, `ceil(x)`, `round(x)` |
| **Time** | `time_now()`, `sleep(ms)` |
| **Filesystem** | `read_file(path)`, `write_file(path, content)` |
| **Arrays** | `len(arr)`, `push(arr, item)`, `pop(arr)`, `first(arr)`, `last(arr)`, `rest(arr)` |

---

## Getting Started

### Prerequisites

- Go 1.18 or higher installed on your system.

### Build

Compile the `rustx` CLI binary:

```bash
go build -o rustx main.go
```

### Running Tests

Execute the unit test suite across all packages:

```bash
go test ./...
```

---

## CLI Usage

### Run a Script

Execute any `.rx` file:

```bash
./rustx run samples/hello.rx
# or shorthand:
./rustx samples/hello.rx
```

### Interactive REPL

Launch the REPL prompt:

```bash
./rustx repl
# or simply:
./rustx
```

Example session:

```text
RustX Starter v0.1
Type 'exit' to quit.
>>> x := 5
>>> y := 10
>>> println("Sum is:", x + y)
Sum is: 15
>>> func greet(name string) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### Check Syntax (Build Check)

Validate syntax without execution:

```bash
./rustx build samples/hello.rx
```

### Version

```bash
./rustx version
```

---

## Future Extensibility

RustX Starter is designed as a foundational stepping stone:

1. **Borrow Checker & Ownership**: The AST maintains token coordinates and explicit statement hierarchies, allowing ownership and lifetime analysis passes to precede evaluation/compilation.
2. **Bytecode VM / Native Compilation**: The clean separation between the AST (`/ast`), parser (`/parser`), and evaluation model (`/object`, `/evaluator`) enables substituting a bytecode compiler or LLVM code generator without touching grammar or parsing logic.
3. **Structs, Traits & Generics**: The recursive descent Pratt parser is ready to accept `type StructName struct { ... }` and `trait` definitions by adding dedicated prefix and statement parsers.
