# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go, implemented in Go. It is designed to be a lightweight scripting + systems hybrid language foundation that can be extended with features like ownership and borrow checking in the future.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array` (internal).
- **Operators**: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `&&`, `||`, `++`.
- **Functions**: Named and anonymous functions, parameters, return types, and recursion.
- **Control Flow**: `if`/`else`, `while`, and `for` loops (Go-style).
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Standard Library**:
    - **Math**: `sqrt`, `abs`, `pow`, `sin`, `cos`.
    - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`.
    - **Time**: `now`, `sleep`.
    - **Filesystem**: `read_file`, `write_file`.

## Architecture

The project follows a modular design:
- `/token`: Token definitions and lookup.
- `/lexer`: Tokenizes the source code.
- `/ast`: Abstract Syntax Tree node definitions.
- `/parser`: Pratt-style recursive descent parser.
- `/object`: Object-based runtime system and environment.
- `/evaluator`: Tree-walking interpreter.
- `/repl`: Interactive REPL with multiline support.
- `main.go`: CLI entry point.

## Usage

### Run a Script
```bash
go run main.go run examples/factorial.rx
```

### Start REPL
```bash
go run main.go repl
# or just
go run main.go
```

### Version
```bash
go run main.go version
```

## Examples

### Factorial
```rustx
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### Loops
```rustx
for i := 0; i < 5; i++ {
    println(i)
}
```

### Filesystem
```rustx
write_file("hello.txt", "Hello RustX")
content := read_file("hello.txt")
println(content)
```
