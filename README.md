# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

RustX Starter uses the `.rx` file extension and runs either as interpreted scripts or via an interactive REPL.

## Key Features

1. **Variables & Mutability**: Mutable by default, supports explicit `var name = value` declaration or Go-style shorthand `name := value`.
2. **Standard Data Types**: Supports `int`, `float`, `bool`, `string`, `Array` lists, and zero-based bracket indexing `[index]` for strings and arrays.
3. **Advanced Control Flow**:
   - Standard unparenthesized `if`/`else` branches.
   - Standard unparenthesized `while` loops.
   - Go-style 3-part `for` loops (e.g. `for i := 0; i < 10; i++`).
4. **First-Class Functions**: Supports optional parameter/return type signatures, anonymous function literals, named functions, and full lexical scope closures and recursion.
5. **Standard Library Built-ins**:
   - **Printing**: `print`, `println`
   - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
   - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
   - **Time**: `sleep`, `now`
   - **Filesystem**: `read_file`, `write_file`
6. **Comments**: Supports single-line comments (`//`) and multi-line/block comments (`/* */`).
7. **Robust Multiline REPL**: Preserves state/variables between executions, balances open braces `{` and parentheses `(` dynamically before evaluating multiline input blocks, and prints error logs gracefully.

## Architecture

The project has a highly clean, modular, and standard compiler/interpreter architecture:

- `/token`: Defines Token types, position information (line & column), and keywords.
- `/lexer`: Implements the tokenizer with full position-aware scanning and comment filtering.
- `/ast`: Defines the position-aware Abstract Syntax Tree nodes for recursive parsing and evaluation.
- `/parser`: Implements a hybrid Recursive Descent & Pratt parser. Configured to stop parsing expressions at `{` to allow clean, unparenthesized control flows.
- `/object`: Defines the runtime object system and the lexically-enclosed `Environment` which recursively updates variables across parent scopes on assignment.
- `/evaluator`: Tree-walking interpreter that implements logical short-circuiting, safe zero-division handling, standard library registration, and expression evaluation.
- `/repl`: State-preserving multiline-capable interactive shell.
- `/samples`: Sample scripts demonstrating the various language features.
- `main.go`: The CLI entry point supporting `run`, `repl`, `version`, and direct shorthand script execution.

---

## Build & Install Instructions

Ensure you have [Go](https://go.dev/) (v1.16+) installed.

To compile the `rustx` executable:
```bash
go build -o rustx main.go
```

To run all unit tests:
```bash
go test ./...
```

---

## Usage

### 1. Shorthand Direct File Execution
```bash
./rustx samples/hello.rx
```

### 2. Using `run` Command
```bash
./rustx run samples/factorial.rx
```

### 3. Interactive REPL
```bash
./rustx repl
```
Example session:
```text
RustX Starter v0.1
>>> x := 5
>>> println(x)
5
>>> func fact(n int) int {
...     if n <= 1 { return 1; }
...     return n * fact(n - 1);
... }
>>> println(fact(x))
120
>>> exit
```

---

## Example Programs

All examples can be found under the `/samples` directory:
- `hello.rx`: Basic hello world script.
- `factorial.rx`: Recursive factorial function example.
- `variables.rx`: Variable scopes, definitions, and mutations.
- `loops.rx`: Demonstration of `for` and `while` loop syntax.
- `std_lib_demo.rx`: Demonstrates string, math, time, and filesystem standard library integrations.
