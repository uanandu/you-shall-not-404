.PHONY: build clean

WASM_OUT  := ../../public/shell.wasm
EXEC_OUT  := ../../public/wasm_exec.js
GOROOT    := $(shell go env GOROOT)

build:
	@echo "→ Compiling Go → WebAssembly..."
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) .
	@echo "→ Copying wasm_exec.js runtime..."
	@cp "$(GOROOT)/misc/wasm/wasm_exec.js" "$(EXEC_OUT)" 2>/dev/null || \
	 cp "$(GOROOT)/lib/wasm/wasm_exec.js"  "$(EXEC_OUT)" 2>/dev/null || \
	 (echo "ERROR: cannot find wasm_exec.js in Go installation" && exit 1)
	@echo "✓ Done. Output: public/shell.wasm + public/wasm_exec.js"
	@ls -lh $(WASM_OUT)

clean:
	rm -f $(WASM_OUT) $(EXEC_OUT)
