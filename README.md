# RustX Starter

**RustX Starter** is a minimal, clean, and scalable programming language inspired by Rust and Go. It combines the clean expression syntax of Rust with the straightforward imperative keywords of Go (`func`, `var`, `:=`).

## Architecture Overview

The system is designed with a modular architecture to allow easy future extensions (such as ownership systems, borrow checkers, bytecode compilation, or async execution) without requiring a rewrite of the parser or AST:

- `token/`: Token definitions, keyword maps, and source location metadata (`Line`, `Column`).
- `errors/`: Standardized error diagnostic types (`SyntaxError`, `RuntimeError`).
- `lexer/`: Tokenizer supporting numbers (int/float), strings with escape sequences (`\n`, `\t`, `\"`, etc.), multi-character operators (`:=`, `==`, `!=`, `<=`, `>=`, `&&`, `||`, `++`, `--`), single-line (`//`), and multi-line (`/* */`) comments.
- `ast/`: AST interface definitions and node implementations supporting statements, expressions, optional type annotations, and future extensions.
- `parser/`: Recursive descent Pratt parser supporting operator precedence, optional semicolons, unparenthesized control flow (`if`, `while`, `for`), statement initializers in for loops, and line/column error tracking.
- `object/`: Runtime object system and block-scoped `Environment` supporting nested variable resolution.
- `evaluator/`: AST tree-walking interpreter and standard built-in functions.
- `repl/`: Interactive multiline REPL preserving state across statements.
- `main.go`: CLI entry point providing `run`, `repl`, `version`, `build` commands, and direct file execution shorthand (`rustx file.rx`).

## Installation & Build Instructions

Requires Go 1.20+.

```bash
# Build the rustx binary
go build -o rustx main.go
```

## CLI Usage

```bash
# Start REPL
./rustx
./rustx repl

# Run a script
./rustx run samples/hello.rx
./rustx samples/hello.rx

# Check syntax (build)
./rustx build samples/hello.rx

# View version
./rustx version
```

## Language Features & Examples

### Hello World (`samples/hello.rx`)
```rust
func main() {
    println("Hello, RustX Starter!")
}
```

### Recursion & Functions (`samples/factorial.rx`)
```rust
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

### Control Flow (`samples/control_flow.rx`)
```rust
func main() {
    var age = 16
    if age >= 18 {
        println("Adult")
    } else {
        println("Minor")
    }

    for i := 0; i < 5; i++ {
        println("i =", i)
    }
}
```

### Built-in Library Functions
- **Printing**: `print(...)`, `println(...)`
- **String**: `len(str)`, `contains(str, substr)`, `split(str, sep)`, `join(array, sep)`, `to_upper(str)`, `to_lower(str)`
- **Math**: `abs(val)`, `pow(base, exp)`, `sin(x)`, `cos(x)`, `tan(x)`
- **Time**: `now()`, `sleep(ms)`
- **Filesystem**: `read_file(path)`, `write_file(path, content)`

## Future Extensibility

RustX Starter's AST nodes, token pipeline, and environment scope layers are structured to seamlessly accommodate future language enhancements:
- Ownership & Borrow Checker
- Async / Await and Goroutine-style Concurrency
- Structs, Enums, and Traits/Interfaces
- Module / Import system
- Native compiler / Bytecode VM backends
