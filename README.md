# RustX Starter

RustX Starter is a compiled/interpreted hybrid programming language designed as a minimal, clean, and scalable foundation inspired by Rust and Go. It balances simplicity, speed, readability, and future extensibility.

---

## Language Features

- **Blended Syntax**: Combines Go and Rust syntax constructs (`func`, `var`, `:=`, optional semicolons, curly brace blocks).
- **Type Annotations**: Optional strong typing in function parameters and return signatures (`func add(a int, b int) int`).
- **Data Types**: `int`, `float`, `string`, `bool`, `null`, and `array`.
- **Control Flow**: `if` / `else if` / `else`, `while`, and 3-part or condition-only `for` loops.
- **Functions**: First-class functions, recursion support, and automatic entrypoint invocation via zero-argument `main()`.
- **Standard Library**:
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `now`, `sleep`
  - Filesystem: `read_file`, `write_file`
- **Error Handling**: Friendly human-readable `SyntaxError` and `RuntimeError` messages with accurate line and column positions.

---

## Architecture Overview

RustX Starter follows a modular architectural pipeline:

```
Source Code (.rx) -> Lexer -> Tokens -> Parser (Pratt) -> AST -> Evaluator / Runtime -> Output
```

1. `token/`: Defines lexical tokens and keywords.
2. `errors/`: Formats human-readable syntax and runtime errors with source coordinates.
3. `lexer/`: Converts source text into a stream of tokens, preserving line and column location info.
4. `ast/`: Abstract Syntax Tree definitions designed for future extensions (ownership, borrow checking, generics, async).
5. `parser/`: Recursive descent parser using Pratt parsing for expressions and operator precedence.
6. `object/`: Object representation system and lexical Environment scope manager.
7. `evaluator/`: Interpreter core, short-circuiting logic, standard library built-in bindings, and entrypoint handler.
8. `repl/`: Interactive shell supporting multiline input buffering and persistent context.

---

## Build & Run Instructions

### Prerequisites
- Go 1.20+ installed

### Building the Executable

```bash
go build -o rustx main.go
```

### CLI Commands

#### 1. Execute a Script
```bash
./rustx run samples/hello.rx
# Or directly via shorthand:
./rustx samples/hello.rx
```

#### 2. Start REPL Prompt
```bash
./rustx repl
# Or with no arguments:
./rustx
```

#### 3. Syntax Verification (Build Check)
```bash
./rustx build samples/factorial.rx
```

#### 4. Version Check
```bash
./rustx version
```

---

## Running Tests

Run all unit test suites across lexer, parser, and evaluator packages:

```bash
go test ./...
```

---

## Example Programs

### Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### Recursive Factorial (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println(factorial(5))
}
```

---

## Future Extensibility

RustX Starter is designed to serve as the foundation for future language features:
- Ownership & Borrow Checker
- Async/Await & Concurrency
- Structs, Enums, and Traits
- Bytecode Virtual Machine / Native Compiler target
