# WASM Shell — Technical Reference

An interactive terminal embedded in the portfolio site. The shell engine is written in Go, compiled to WebAssembly, and driven by a Svelte 5 UI component. Visitors can explore a curated virtual filesystem containing resume content, project documentation, and system information.

---

## 1. Overview

**Why Go + WASM?**

- Go's `syscall/js` package provides a clean, low-ceremony bridge between Go and the browser's JS environment.
- The resulting binary is self-contained: no npm packages, no runtime dependencies beyond the Go WASM shim (`wasm_exec.js`).
- Go compiles to a single `.wasm` file. The shell engine (VFS, command dispatch, pipe support) runs natively in the browser at near-native speed after load.
- The approach produces a real shell with stateful cwd, pipes, tab completion, and history — not a switch/case command lookup pretending to be a terminal.

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────┐
│  Browser                                            │
│                                                     │
│  ┌──────────────────────┐                           │
│  │  Terminal.svelte      │  Svelte 5 island         │
│  │  (src/components/     │  Renders lines, input    │
│  │   svelte/)            │  Handles keyboard events │
│  └──────────┬───────────┘                           │
│             │  window.wasmRunCommand(cmd)            │
│             │  window.wasmGetCompletions(partial)    │
│             ▼                                       │
│  ┌──────────────────────┐                           │
│  │  JS ↔ WASM Bridge    │  syscall/js exports       │
│  │  (main.go)           │  Registered on window     │
│  └──────────┬───────────┘                           │
│             │                                       │
│             ▼                                       │
│  ┌──────────────────────┐                           │
│  │  Shell Engine        │  shell.go                 │
│  │  · cwd state         │  Tokenisation, pipes,     │
│  │  · env map           │  history, completions     │
│  │  · history slice     │                           │
│  └──────────┬───────────┘                           │
│             │                                       │
│      ┌──────┴──────┐                               │
│      ▼             ▼                               │
│  ┌────────┐  ┌──────────┐                          │
│  │  VFS   │  │ Commands │  fs.go / commands.go      │
│  │ fs.go  │  │ map      │  seed.go                  │
│  └────────┘  └──────────┘                          │
│                                                     │
│  Static assets (public/)                           │
│    shell.wasm       ← compiled Go binary            │
│    wasm_exec.js     ← Go's JS runtime shim          │
└─────────────────────────────────────────────────────┘
```

**Data flow for a single command:**

1. User presses Enter in Terminal.svelte.
2. Svelte calls `window.wasmRunCommand(input)`.
3. Go's `js.FuncOf` handler receives the string, calls `sh.run(input)`.
4. `sh.run` tokenises, resolves pipes, dispatches to the `commands` map or built-ins.
5. The command reads/walks the VFS, formats output, returns a `result` struct.
6. `resultObj` builds a plain JS object `{output, err, clear}` and returns it.
7. Svelte reads the object and appends output lines to its reactive `lines` array.

---

## 3. Virtual Filesystem

The VFS is a pure in-memory tree. There is no IndexedDB, no localStorage, no real filesystem access. It is read-only by design.

### Node structure (`fs.go`)

```go
type node struct {
    name    string
    kind    nodeKind      // kindDir | kindFile
    content string        // file body (empty for dirs)
    mtime   time.Time
    mode    string        // e.g. "-rw-r--r--", "drwxr-xr-x"
    entries map[string]*node  // nil for files; child nodes for dirs
}
```

Directories are created with `dir(name, children...)` and files with `file(name, content)`. Both are variadic helper functions in `fs.go`. The `file()` helper strips a leading newline from the content string, allowing the content literal to start on the line after the backtick without producing a blank first line.

### Path resolution (`vfs.resolve`)

```
vfs.resolve(cwd, path) → (*node, absPath, error)
```

- If `path` is empty or `~`, resolves to `/home/anandu`.
- If `path` starts with `~/`, expands to `/home/anandu/<rest>`.
- If `path` is not absolute, prepends `cwd`.
- Calls `path.Clean` to normalise `..` and `.` components.
- Walks the tree segment by segment, returning an error if any segment is missing.

### Filesystem layout

```
/
├── home/
│   └── anandu/
│       ├── README.md
│       ├── skills.txt
│       ├── contact.txt
│       ├── projects/
│       │   ├── README.md
│       │   ├── api-darwin/
│       │   │   ├── README.md
│       │   │   └── main.go
│       │   ├── devrouter/
│       │   │   ├── README.md
│       │   │   └── devrouter.yml
│       │   └── infra-mcp/
│       │       └── README.md
│       └── experience/
│           ├── README.md
│           ├── deloitte.txt
│           ├── off-paper.txt
│           ├── ypl.txt
│           └── techni-mek.txt
├── etc/
│   ├── os-release
│   └── hostname
└── usr/
    └── local/
        └── bin/
            ├── api-darwin
            └── devrouter
```

---

## 4. Shell Engine (`shell.go`)

### State

```go
type shell struct {
    fs      *vfs
    cwd     string
    env     map[string]string
    history []string
}
```

`cwd` starts at `/home/anandu`. The `cd` command mutates it in place.

### Command parsing

`tokenise(s string) []string` splits by whitespace, respecting double-quoted strings. A quoted string like `"hello world"` is kept as a single token.

### Pipe support

`splitPipe(s string) []string` splits the input on `|` characters (not inside quotes). Each stage is executed in sequence; the `output` string from stage N becomes the `stdin` argument to stage N+1. If any stage returns `isErr: true`, the pipeline aborts and the error is returned immediately.

### Built-ins

`cd`, `history`, `clear`, `exit`, and `logout` are handled directly in `exec()` rather than through the `commands` map, because:
- `cd` must mutate `sh.cwd`.
- `history` must read `sh.history`.
- `clear` returns `result{clear: true}`, which tells the Svelte UI to wipe the line buffer.

### Tab completion (`sh.completions`)

- If the partial input has no spaces, completes command names from the `commands` map.
- If the partial has spaces, completes the last token as a filesystem path. The parent directory is resolved, and all entries whose names start with the typed prefix are returned. Directories get a trailing `/`.

---

## 5. Commands

| Command | Synopsis | Notes |
|---------|----------|-------|
| `help` | `help` | Lists all commands |
| `ls` | `ls [-la] [path]` | `-l` long format, `-a` include dotfiles |
| `cd` | `cd [path]` | Mutates cwd; `~` goes home |
| `pwd` | `pwd` | Prints cwd |
| `cat` | `cat [file...]` | Concatenates and prints files |
| `head` | `head [-n N] [file]` | First N lines (default 10) |
| `grep` | `grep [-i] pattern [file]` | `-i` case-insensitive; works from stdin |
| `wc` | `wc [file]` | Lines, words, chars |
| `find` | `find [path] [-name pattern]` | Recursive tree walk |
| `echo` | `echo [args...]` | Expands `$VAR` from env map |
| `whoami` | `whoami` | Prints `$USER` |
| `env` | `env` | Prints all env vars, sorted |
| `man` | `man <cmd>` | Manual pages for select commands |
| `uname` | `uname [-a]` | OS name / full info |
| `date` | `date` | Browser's current time via `time.Now()` |
| `neofetch` | `neofetch` | ASCII art system overview |
| `history` | `history` | Numbered command history |
| `clear` | `clear` | Clears the terminal (sets `clear: true`) |
| `exit` / `logout` | — | Friendly refusal |
| `__welcome__` | internal | Called by UI on startup |

**Pipe example:**
```
cat skills.txt | grep Go
cat experience/deloitte.txt | wc
ls projects/ | grep darwin
```

---

## 6. JS ↔ WASM Bridge (`main.go`)

### How `syscall/js` exports work

Go functions are registered on the browser's `window` object using `js.Global().Set(name, js.FuncOf(fn))`. Once registered, any JS code can call `window.wasmRunCommand("ls")` directly.

`js.FuncOf` wraps a Go function with the signature:

```go
func(_ js.Value, args []js.Value) any
```

The first argument is the JS `this` value (ignored). `args` is the list of JS arguments passed at the call site.

### The `resultObj` pattern

```go
func resultObj(output string, isErr, clear bool) js.Value {
    obj := js.Global().Get("Object").New()
    obj.Set("output", output)
    obj.Set("err", isErr)
    obj.Set("clear", clear)
    return obj
}
```

This creates a plain JS object (not a Go struct) so the Svelte side can access `.output`, `.err`, `.clear` without any serialisation step. Returning a `js.Value` means no JSON round-trip.

### Why `select {}`

WebAssembly modules must not return from `main()` while JS callbacks are still registered. If `main` exits, the Go runtime tears down and the registered `js.Func` callbacks become invalid, causing panics on any subsequent call. The empty `select{}` blocks the goroutine forever, keeping the WASM module alive for the lifetime of the page.

### Completions return `[]any`

`wasmGetCompletions` returns a JS array. Go's `syscall/js` does not automatically convert `[]string` to a JS array, so each string is boxed into `any` and passed to `js.ValueOf([]any{...})`, which produces a proper JS `Array`.

---

## 7. Build Process

### Prerequisites

- Go 1.22 or later (`go version`)
- `GOROOT` must be set correctly (it is, if Go is installed via `go install` or a package manager)

### Steps

```bash
# 1. Navigate to the shell source
cd wasm/shell

# 2. Build the WASM binary and copy the runtime shim
make build
```

`make build` runs two operations:

1. **Compile:**
   ```
   GOOS=js GOARCH=wasm go build -o ../../public/shell.wasm .
   ```
   Sets the Go cross-compilation targets to `js/wasm`. Outputs `public/shell.wasm`.

2. **Copy shim:**
   ```
   cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../../public/wasm_exec.js
   ```
   (Falls back to `lib/wasm/wasm_exec.js` for newer Go versions.) The shim is Go's standard JS glue that bootstraps the WASM runtime, sets up memory, and implements the syscall interface.

### Output

```
public/
  shell.wasm      ~5–10 MB (typical for a Go WASM binary)
  wasm_exec.js    ~17 KB
```

Both are served as static assets by Astro / Cloudflare Pages.

### Clean

```bash
make clean   # removes public/shell.wasm and public/wasm_exec.js
```

---

## 8. Adding Commands

1. Write the implementation in `commands.go` (or a new file in the same package):

```go
func cmdMyCommand(sh *shell, args []string, stdin string) result {
    // args: tokens after the command name
    // stdin: piped input from previous stage (empty if none)
    if len(args) == 0 {
        return err("mycommand: missing argument")
    }
    return ok("Hello, " + args[0])
}
```

2. Register it in the `commands` map at the top of `commands.go`:

```go
var commands = map[string]cmdFn{
    // ... existing entries ...
    "mycommand": cmdMyCommand,
}
```

3. Optionally add a man page entry in `manPages`:

```go
var manPages = map[string]string{
    // ... existing entries ...
    "mycommand": "MYCOMMAND(1)\n\nNAME\n  mycommand — does a thing\n\nSYNOPSIS\n  mycommand <name>",
}
```

4. Rebuild: `cd wasm/shell && make build`

**Note:** Commands that need to mutate shell state (like `cd`) cannot be `cmdFn` values. Define them as methods on `*shell` and call them directly in `exec()`.

---

## 9. Adding Filesystem Content

All content lives in `seed.go`. The `seedFS()` function returns the root node.

**Add a file:**
```go
file("resume.txt", `
Anandu
Montréal, QC
...
`),
```

**Add a directory with files:**
```go
dir("notes",
    file("ideas.txt", `
- Build something with Rust
- Learn more about eBPF
`),
    file("reading.txt", `
Current: "Crafting Interpreters" — Robert Nystrom
`),
),
```

Place the new `dir()` or `file()` call inside the appropriate parent `dir()` call in `seedFS()`. The `file()` helper strips the leading newline from the content literal.

After editing `seed.go`, rebuild: `cd wasm/shell && make build`

---

## 10. Limitations

This is a **simulated shell**, not a real one. Key non-features:

| What it is not | Why |
|----------------|-----|
| A real process model | There is no `fork`, no `exec`, no subprocess spawning. Commands are Go functions. |
| Networked | WASM cannot open raw sockets. `curl`, `wget`, `ping` are not possible. |
| Persistent | No state survives a page reload. History, cwd, and any mutations are in-memory only. |
| Writable | The VFS is read-only. `touch`, `mkdir`, `rm`, `echo foo > file` are not implemented. |
| POSIX-compliant | Tokenisation is simplified. No glob expansion, no variable assignment, no subshells, no redirects. |
| A real kernel | There is no signal handling beyond `Ctrl+C` (which just clears the input line). |

---

## 11. Performance Notes

### Binary size

A minimal Go WASM binary is typically 5–10 MB. This is larger than a comparable Rust binary (~1–2 MB for similar functionality), because Go includes its runtime, garbage collector, and goroutine scheduler in every binary. For a portfolio terminal this is acceptable — the binary is cached after the first load.

Optimisations that can reduce size:
- `go build -ldflags="-s -w"` strips debug symbols and DWARF info, reducing size by ~20–30%.
- `wasm-opt` (from the Binaryen toolkit) can further reduce size by ~15%.

### Instantiation time

`WebAssembly.instantiateStreaming` begins parsing as the binary streams in, so instantiation starts before the full binary is downloaded. On a fast connection, the shell is ready within 1–2 seconds. On a slow connection, the "Loading shell.wasm..." indicator is shown.

### After load

Once instantiated, all command execution is synchronous Go function calls — no network, no disk. Response is effectively instantaneous for all commands. The WASM runtime's garbage collector is the only variable, but the shell's memory footprint is tiny (a few hundred KB of live objects at most).

### Caching

`shell.wasm` is a static asset. Both Cloudflare Pages and Vercel serve it with a long `Cache-Control` header on subsequent visits. `wasm_exec.js` is also cached. After the first visit, the shell loads from the browser cache in milliseconds.
