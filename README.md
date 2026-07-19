# RustX Starter

RustX Starter is a minimal, clean, and highly extensible programming language inspired by **Rust** and **Go**. It focuses on simplicity, speed, readability, and future extensibility.

Currently, RustX Starter behaves as a lightweight scripting + systems hybrid language implemented in Go. It lays down a robust foundation of parser, lexer, and runtime design, so that advanced features such as ownership, lifetimes, structs, and a bytecode VM can be added in the future.

---

## Features

1. **Variables**:
   - Explicit mutable declaration: `var x = 5`
   - Shorthand declaration: `y := 10`
   - Optional type annotations: `var name string = "Muhammad"` (parsed but ignored during initial tree-walking phase)

2. **Data Types**:
   - `int` (64-bit integer)
   - `float` (64-bit float)
   - `bool` (boolean)
   - `string` (UTF-8 string)
   - `array` (dynamic lists: e.g. `[1, 2, 3]`)

3. **Operators**:
   - Arithmetic: `+`, `-`, `*`, `/`
   - Postfix: `++`, `--` (limited to identifiers)
   - Comparisons: `==`, `!=`, `<`, `>`, `<=`, `>=`
   - Logical with short-circuiting: `&&`, `||`

4. **Functions**:
   - Parameters and return types
   - Self-referencing recursive functions
   - Named and anonymous functions
   - Lexical closures and block scoping

5. **Control Flow**:
   - Unparenthesized / parenthesized `if`, `else if`, `else`
   - Loops: `while` and flexible Go-style `for` loops (`for i := 0; i < 10; i++ { ... }`)

6. **Standard Library**:
   - `print(...)` & `println(...)`
   - strings: `len(...)`, `contains(str, sub)`, `split(str, sep)`, `join(array, sep)`, `to_upper(str)`, `to_lower(str)`
   - math: `abs(...)`, `pow(base, exp)`, `sin(...)`, `cos(...)`, `tan(...)`
   - time: `now()`, `sleep(seconds)`
   - filesystem: `read_file(path)`, `write_file(path, content)`

---

## Directory Structure

- `errors/`: Custom diagnostic Syntax and Runtime error representations with precise line/column locations.
- `token/`: Lexical token definitions, literals, and keyword lookups.
- `lexer/`: Lexer / Tokenizer tracking source coordinates (line/column) and skipping line & block comments.
- `ast/`: AST Node interfaces implementing `GetToken()` for position-aware error mapping.
- `parser/`: Pratt-based recursive descent parser supporting unparenthesized control structure parsing.
- `object/`: Type representations in RustX's runtime environment (closures, environments, arrays, and errors).
- `evaluator/`: Main tree-walking interpreter flow, standard library built-ins registry, short-circuit operators, and safe division-by-zero checks.
- `repl/`: REPL interactive prompt preserving variable environments and managing balanced multiline entry blocks.
- `samples/`: Sample RustX scripts demonstrating the language capabilities.

---

## Getting Started

### Installation
Ensure you have Go installed (Go 1.24+ is recommended). Clone the repository and run:

```bash
go build -o rustx .
```

### Run the CLI Commands
The compiled binary `rustx` provides three main modes:

1. **Start REPL**:
   ```bash
   ./rustx repl
   # or simply
   ./rustx
   ```
   *You can write multi-line blocks inside the REPL and they will evaluate when all braces/parentheses are matched!*

2. **Execute a script file**:
   ```bash
   ./rustx run samples/factorial.rx
   # or use shorthand directly
   ./rustx samples/factorial.rx
   ```

3. **Check Version**:
   ```bash
   ./rustx version
   ```

---

## Error Diagnostics Example

RustX provides highly readable, precise error diagnostics.

For example, a syntactic issue will output:
```text
SyntaxError: expected next token to be RPAREN, got INT instead at line 5 column 12
```

A dynamic issue like division-by-zero or out of bounds will output:
```text
RuntimeError: division by zero at line 3 column 10
```

---

## Future Extensibility

The parser and interpreter architectures have been intentionally isolated:
- AST nodes carry full metadata context so a **static ownership analyzer / borrow checker** can pass over the generated AST before execution.
- The tree-walker can easily be converted to compile AST nodes into a flat bytecode instruction format for execution on a high-performance register-based or stack-based **Bytecode VM**.
