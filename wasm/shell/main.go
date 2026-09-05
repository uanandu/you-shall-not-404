//go:build js && wasm

package main

import (
	"syscall/js"
)

func main() {
	sh := newShell()

	js.Global().Set("wasmRunCommand", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) == 0 {
			return resultObj("", false, false)
		}
		r := sh.run(args[0].String())
		return resultObj(r.output, r.isErr, r.clear)
	}))

	js.Global().Set("wasmGetCompletions", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) == 0 {
			return js.ValueOf([]any{})
		}
		completions := sh.completions(args[0].String())
		items := make([]any, len(completions))
		for i, c := range completions {
			items[i] = c
		}
		return js.ValueOf(items)
	}))

	js.Global().Set("wasmGetCwd", js.FuncOf(func(_ js.Value, _ []js.Value) any {
		return js.ValueOf(sh.cwd)
	}))

	// Block forever — WASM must not exit
	select {}
}

func resultObj(output string, isErr, clear bool) js.Value {
	// Build a plain JS object via js.Global().Get("Object")
	obj := js.Global().Get("Object").New()
	obj.Set("output", output)
	obj.Set("err", isErr)
	obj.Set("clear", clear)
	return obj
}

