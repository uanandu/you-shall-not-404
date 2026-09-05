package main

import (
	"strings"
)

type result struct {
	output string
	isErr  bool
	clear  bool
}

type shell struct {
	fs      *vfs
	cwd     string
	env     map[string]string
	history []string
}

func newShell() *shell {
	sh := &shell{
		fs:  &vfs{root: seedFS()},
		cwd: "/home/anandu",
		env: map[string]string{
			"HOME":  "/home/anandu",
			"USER":  "anandu",
			"SHELL": "/bin/gosh",
			"TERM":  "xterm-256color",
			"LANG":  "en_CA.UTF-8",
			"PATH":  "/usr/local/bin:/usr/bin:/bin",
		},
	}
	return sh
}

func (sh *shell) run(input string) result {
	input = strings.TrimSpace(input)
	if input == "" {
		return result{}
	}

	// Record history (skip duplicates at head)
	if len(sh.history) == 0 || sh.history[len(sh.history)-1] != input {
		sh.history = append(sh.history, input)
	}

	// Handle pipes: split by | and chain
	stages := splitPipe(input)
	if len(stages) == 1 {
		return sh.exec(stages[0], "")
	}

	// Simple pipe: stdout of stage N is stdin of stage N+1
	var buf string
	for _, stage := range stages {
		r := sh.exec(strings.TrimSpace(stage), buf)
		if r.isErr {
			return r
		}
		buf = r.output
	}
	return result{output: buf}
}

func (sh *shell) exec(input, stdin string) result {
	tokens := tokenise(input)
	if len(tokens) == 0 {
		return result{}
	}

	cmd, args := tokens[0], tokens[1:]

	// Built-in: cd must mutate shell state
	if cmd == "cd" {
		return sh.cmdCD(args)
	}
	if cmd == "history" {
		return sh.cmdHistory()
	}
	if cmd == "clear" {
		return result{clear: true}
	}
	if cmd == "exit" || cmd == "logout" {
		return result{output: "Use the browser tab to close. 👋"}
	}

	fn, ok := commands[cmd]
	if !ok {
		// Check if it looks like a path execution
		if strings.HasPrefix(cmd, "./") || strings.HasPrefix(cmd, "/") {
			return sh.cmdExec(cmd, args)
		}
		return result{output: cmd + ": command not found", isErr: true}
	}
	return fn(sh, args, stdin)
}

func (sh *shell) completions(partial string) []string {
	// Complete command names if no space, else complete paths
	if !strings.Contains(partial, " ") {
		var out []string
		for name := range commands {
			if strings.HasPrefix(name, partial) {
				out = append(out, name)
			}
		}
		return out
	}

	// Path completion: last token
	tokens := tokenise(partial)
	last := ""
	if len(tokens) > 0 {
		last = tokens[len(tokens)-1]
	}

	// Find parent dir
	parentPath := sh.cwd
	prefix := ""
	if idx := strings.LastIndex(last, "/"); idx >= 0 {
		parentPath = sh.fs.abs(sh.cwd, last[:idx+1])
		prefix = last[:idx+1]
	}

	parent, err := sh.fs.walk(parentPath)
	if err != nil || parent.kind != kindDir {
		return nil
	}

	var out []string
	for name, child := range parent.entries {
		if strings.HasPrefix(name, strings.TrimPrefix(last, prefix)) {
			suffix := ""
			if child.kind == kindDir {
				suffix = "/"
			}
			out = append(out, prefix+name+suffix)
		}
	}
	return out
}

// tokenise splits by spaces, handling basic double-quoted strings
func tokenise(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	for _, r := range s {
		switch {
		case r == '"' && !inQuote:
			inQuote = true
		case r == '"' && inQuote:
			inQuote = false
		case r == ' ' && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

func splitPipe(s string) []string {
	// Split by | but not inside quotes
	var stages []string
	var cur strings.Builder
	inQuote := false
	for _, r := range s {
		switch {
		case r == '"':
			inQuote = !inQuote
			cur.WriteRune(r)
		case r == '|' && !inQuote:
			stages = append(stages, cur.String())
			cur.Reset()
		default:
			cur.WriteRune(r)
		}
	}
	stages = append(stages, cur.String())
	return stages
}
