# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable interpreted/scripting programming language inspired by the elegance of **Rust** and **Go**. It blends high readability, explicit control, and performance with a modular architecture designed for future extensibility (such as memory safety systems, borrow checking, and static typing extensions).

---

## Architecture & Modular Design

RustX Starter follows standard compiler/interpreter design principles divided into clean, decoupled packages:

```
rustx/
├── token/       # Token type definitions, keywords, operators, position tracking
├── errors/      # Standardized human-readable SyntaxError and RuntimeError diagnostics
├── lexer/       # Lexical analyzer with comment skipping and escape sequence support
├── ast/         # Abstract Syntax Tree nodes for expressions and statements
├── parser/      # Pratt recursive-descent parser supporting optional semicolons
├── object/      # Runtime value representation (Integer, Float, String, Function, etc.)
├── evaluator/   # AST Tree-walking interpreter & Built-in standard library functions
├── repl/        # Interactive shell supporting multiline input and state preservation
├── samples/     # Example programs demonstrating recursion, loops, and file I/O
└── main.go      # CLI entry point (run, repl, build, version)
```

---

## Key Features & Syntax

- **Variable Declarations:**
  - `var` declaration: `var age int = 16` or `var age = 16`
  - Short-hand assignment: `name := "Muhammad"`
  - Direct assignment: `count = count + 1`

- **Data Types:**
  - `int` (64-bit integer)
  - `float` (64-bit floating point)
  - `bool` (`true`, `false`)
  - `string` (Double-quoted with escape sequences like `\n`, `\t`, `\"`)
  - `array` (`[1, 2, 3]`)

- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`
  - Postfix Increment/Decrement: `++`, `--`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`

- **Functions & Control Flow:**
  - Functions: `func name(arg type) return_type { ... }`
  - Recursion support
  - Loops: `for i := 0; i < 10; i++ { ... }` or `while condition { ... }`
  - Conditionals: `if`, `else if`, `else`

- **Standard Library Built-ins:**
  - Printing: `print(...)`, `println(...)`
  - Strings: `len(s)`, `contains(str, substr)`, `split(str, sep)`, `join(array, sep)`, `to_upper(s)`, `to_lower(s)`
  - Math: `abs(n)`, `pow(b, e)`, `sin(rad)`, `cos(rad)`, `tan(rad)`
  - Time & System: `now()`, `sleep(ms)`
  - Filesystem: `read_file(path)`, `write_file(path, contents)`

---

## Building and Running

### Build CLI Binary

```bash
go build -o rustx main.go
```

### Running Scripts

Execute a RustX `.rx` script directly:

```bash
./rustx run samples/hello.rx
```

Or using shorthand file execution:

```bash
./rustx samples/factorial.rx
```

### Syntactic Validation (Build command)

```bash
./rustx build samples/loops_and_math.rx
```

### Interactive REPL Prompt

```bash
./rustx repl
```

Example REPL interaction:
```text
RustX Starter v0.1 REPL
Type 'exit' or Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum =", x + y)
Sum = 30
```

---

## Testing

Run all package unit tests across the codebase:

```bash
go test ./...
```

---

## Future Architecture Philosophy

RustX Starter intentionally defers borrow checking, ownership models, and explicit lifetime annotations in this starter phase. However, the parser and AST are structured with extensible node interfaces, making it simple to attach static analysis passes, borrow checking rules, bytecode compilation, or LLVM/C code generation in future iterations without rewriting the core language pipeline.
