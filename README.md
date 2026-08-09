# RustX Starter

RustX Starter is a minimal, clean, and extensible programming language inspired by Go and Rust, designed as a foundation for a modern hybrid scripting and systems language. It features optional static types, flexible assignments, Go-style unparenthesized control structures, nested scopes, recursive functions, arrays, indexing, and a built-in standard library.

The architecture is designed to be highly modular so that advanced systems (such as a borrow checker, concurrency models, or a bytecode VM) can be introduced later without requiring a full compiler or interpreter rewrite.

---

## Language Features

### 1. Variables and Assignments
- Variable declarations with `var` or shorthand declaration with `:=`.
- Block scoping is fully respected: variable lookup falls back to enclosing blocks, and updates via standard `=` recursively trace scope up to parent environments.
```rustx
var x int = 5;
y := "RustX";
x = 10;
```

### 2. Primitive & Composite Types
- Types supported: `int`, `float`, `string`, `bool`.
- Array type `[1, 2, 3]` and zero-based indexing `array[index]` for arrays and strings.

### 3. Control Flow
- Semicolons are optional.
- Parentheses around control flow conditions are optional.
- Supports `if-else`, `while`, and robust `for` loops (including Go-style three-part for loops, single conditional loops, and infinite loops).
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

### 4. Functions & Recursion
- Functions can be named or anonymous, supporting recursion, positional parameters, and optional type annotations.
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
```

### 5. Built-in Standard Library
- **Printing:** `print()`, `println()` (space-separated, newline-terminated output).
- **Strings:** `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`.
- **Math:** `abs()`, `pow()`, `sin()`, `cos()`, `tan()`.
- **Time:** `sleep()` (in milliseconds), `now()` (epoch milliseconds).
- **Filesystem:** `read_file()`, `write_file()`.

---

## Project Architecture

The interpreter is implemented from scratch in Go and is structured into cleanly separated modules:
- `/token`: Defines language tokens, keyword mapping, and location-aware positioning (Line/Column).
- `/errors`: Standardizes user-friendly compiler diagnostics (`SyntaxError` and `RuntimeError`).
- `/lexer`: A stateful tokenizer supporting multi-character symbols, backslash-escaped strings, single-line/multi-line comments, and positional tracking.
- `/ast`: Generates a positional-aware AST (Abstract Syntax Tree) extensible for future compilation targets.
- `/parser`: Uses a top-down Recursive Descent parsing engine combined with a Pratt operator precedence parser. It implements unparenthesized control structure parsing and custom lookahead for for-loop headers.
- `/object`: Our runtime-evaluation system utilizing an Object interface and nested Environment scoping.
- `/evaluator`: Implements a tree-walking interpreter to execute AST nodes, safely executing division-by-zero checks, short-circuiting logical operations, block-scoped loops, and function calls.
- `/repl`: A multiline-aware interactive prompt that tracks block delimiters `{` and `(` before triggering evaluation.

---

## Build & Run Instructions

Ensure Go (1.18+) is installed.

### Build Executable
Build the single binary executable:
```bash
go build -o rustx main.go
```

### Run Scripts
Execute any `.rx` file directly:
```bash
./rustx run samples/hello.rx
# OR shorthand:
./rustx samples/hello.rx
```

### Static Syntax Check
To compile and perform lexical and syntax validation on a file without running it:
```bash
./rustx build samples/factorial.rx
```
If errors are found, the build process exits with status code `1` and prints clean diagnostic messages:
```
SyntaxError:
expected next token to be ), got ( instead
at line 2 column 12
```

### Start the REPL
Launch the multiline interactive prompt:
```bash
./rustx repl
# OR simply run:
./rustx
```
Inside the REPL, blocks can be typed across multiple lines:
```
>>> func greeting() {
...   println("Hello from REPL")
... }
>>> greeting()
Hello from REPL
```

---

## Testing
Unit tests are fully implemented across all processing stages. Run them using:
```bash
go test ./...
```
