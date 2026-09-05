package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type cmdFn func(sh *shell, args []string, stdin string) result

var commands = map[string]cmdFn{
	"help":       cmdHelp,
	"ls":         cmdLS,
	"cat":        cmdCat,
	"pwd":        cmdPwd,
	"echo":       cmdEcho,
	"whoami":     cmdWhoami,
	"env":        cmdEnv,
	"grep":       cmdGrep,
	"man":        cmdMan,
	"uname":      cmdUname,
	"date":       cmdDate,
	"neofetch":   cmdNeofetch,
	"wc":         cmdWC,
	"head":       cmdHead,
	"find":       cmdFind,
	"__welcome__": cmdWelcome,
}

func ok(s string) result  { return result{output: s} }
func err(s string) result { return result{output: s, isErr: true} }

// ── __welcome__ ────────────────────────────────────────────────
func cmdWelcome(_ *shell, _ []string, _ string) result {
	return ok(`
    /\      anandu@portfolio
   /  \     ─────────────────────────────────
  / /\ \    OS       PortfolioOS (WASM/Go)
 / ____ \   Shell    gosh 1.0.0
/_/    \_\  Runtime  Go → WebAssembly
            Role     Fullstack & DevOps Engineer
            Stack    Go · Rust · TypeScript · K8s
            Current  Consultant @ Deloitte Canada
            Web      anandu.dev

  ls · cat README.md · cat skills.txt · help
`)
}

// ── help ───────────────────────────────────────────────────────
func cmdHelp(_ *shell, _ []string, _ string) result {
	return ok(`Available commands:

  Navigation
    ls [-la] [path]     list directory contents
    cd [path]           change directory  (~ = home)
    pwd                 print working directory
    find [path] [name]  find files

  Files
    cat [file...]       display file contents
    head [-n N] [file]  display first N lines (default 10)
    grep [-i] pat file  search file for pattern
    wc [file]           word / line count

  System
    whoami              current user
    uname [-a]          system information
    date                current date and time
    env                 environment variables
    neofetch            system overview (try this)
    history             command history
    clear               clear the terminal

  Meta
    help                this message
    man <command>       command manual
    exit                (does nothing — you're stuck with me)

Tip: tab-complete paths and commands. Pipes work: cat file | grep foo`)
}

// ── ls ─────────────────────────────────────────────────────────
func cmdLS(sh *shell, args []string, _ string) result {
	long := false
	all := false
	target := sh.cwd

	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			if strings.Contains(a, "l") {
				long = true
			}
			if strings.Contains(a, "a") {
				all = true
			}
		} else {
			target = a
		}
	}

	n, abs, e := sh.fs.resolve(sh.cwd, target)
	if e != nil {
		return err(e.Error())
	}
	if n.kind == kindFile {
		if long {
			return ok(fmt.Sprintf("%s  %6d  %s  %s", n.mode, utf8.RuneCountInString(n.content), n.mtime.Format("Jan  2"), n.name))
		}
		return ok(n.name)
	}

	_ = abs
	entries := sh.fs.sortedEntries(n)

	if long {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("total %d\n", len(entries)))
		for _, c := range entries {
			if !all && strings.HasPrefix(c.name, ".") {
				continue
			}
			size := utf8.RuneCountInString(c.content)
			if c.kind == kindDir {
				size = len(c.entries)
			}
			b.WriteString(fmt.Sprintf("%s  %6d  %s  %s\n", c.mode, size, c.mtime.Format("Jan  2"), c.name))
		}
		return ok(strings.TrimRight(b.String(), "\n"))
	}

	// Short listing
	var names []string
	for _, c := range entries {
		if !all && strings.HasPrefix(c.name, ".") {
			continue
		}
		name := c.name
		if c.kind == kindDir {
			name += "/"
		}
		names = append(names, name)
	}
	return ok(strings.Join(names, "  "))
}

// ── cd ─────────────────────────────────────────────────────────
func (sh *shell) cmdCD(args []string) result {
	target := "/home/anandu"
	if len(args) > 0 {
		target = args[0]
	}

	n, abs, e := sh.fs.resolve(sh.cwd, target)
	if e != nil {
		return err("cd: " + e.Error())
	}
	if n.kind != kindDir {
		return err("cd: not a directory: " + target)
	}
	sh.cwd = abs
	return result{}
}

// ── cat ────────────────────────────────────────────────────────
func cmdCat(sh *shell, args []string, stdin string) result {
	if len(args) == 0 {
		if stdin != "" {
			return ok(stdin)
		}
		return err("cat: missing operand")
	}
	var parts []string
	for _, a := range args {
		n, _, e := sh.fs.resolve(sh.cwd, a)
		if e != nil {
			return err("cat: " + e.Error())
		}
		if n.kind == kindDir {
			return err("cat: " + a + ": is a directory")
		}
		parts = append(parts, n.content)
	}
	return ok(strings.Join(parts, "\n"))
}

// ── head ───────────────────────────────────────────────────────
func cmdHead(sh *shell, args []string, stdin string) result {
	n := 10
	var files []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-n" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &n)
			i++
		} else {
			files = append(files, args[i])
		}
	}

	getText := func(p string) (string, bool) {
		nd, _, e := sh.fs.resolve(sh.cwd, p)
		if e != nil {
			return "", false
		}
		return nd.content, true
	}

	var content string
	if len(files) == 0 {
		content = stdin
	} else {
		c, ok2 := getText(files[0])
		if !ok2 {
			return err("head: cannot open '" + files[0] + "'")
		}
		content = c
	}

	lines := strings.Split(content, "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return ok(strings.Join(lines[:n], "\n"))
}

// ── pwd ────────────────────────────────────────────────────────
func cmdPwd(sh *shell, _ []string, _ string) result {
	return ok(sh.cwd)
}

// ── echo ───────────────────────────────────────────────────────
func cmdEcho(sh *shell, args []string, _ string) result {
	// Expand $VAR
	var out []string
	for _, a := range args {
		if strings.HasPrefix(a, "$") {
			key := strings.TrimPrefix(a, "$")
			if v, exists := sh.env[key]; exists {
				out = append(out, v)
				continue
			}
		}
		out = append(out, a)
	}
	return ok(strings.Join(out, " "))
}

// ── whoami ─────────────────────────────────────────────────────
func cmdWhoami(sh *shell, _ []string, _ string) result {
	return ok(sh.env["USER"])
}

// ── env ────────────────────────────────────────────────────────
func cmdEnv(sh *shell, _ []string, _ string) result {
	keys := make([]string, 0, len(sh.env))
	for k := range sh.env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var lines []string
	for _, k := range keys {
		lines = append(lines, k+"="+sh.env[k])
	}
	return ok(strings.Join(lines, "\n"))
}

// ── grep ───────────────────────────────────────────────────────
func cmdGrep(sh *shell, args []string, stdin string) result {
	if len(args) == 0 {
		return err("grep: missing pattern")
	}

	caseInsensitive := false
	var pattern string
	var files []string

	for _, a := range args {
		if a == "-i" {
			caseInsensitive = true
		} else if pattern == "" {
			pattern = a
		} else {
			files = append(files, a)
		}
	}

	search := func(content, src string) []string {
		var matches []string
		matchPat := pattern
		if caseInsensitive {
			matchPat = strings.ToLower(pattern)
		}
		for i, line := range strings.Split(content, "\n") {
			check := line
			if caseInsensitive {
				check = strings.ToLower(line)
			}
			if strings.Contains(check, matchPat) {
				if src != "" {
					matches = append(matches, fmt.Sprintf("%s:%d:%s", src, i+1, line))
				} else {
					matches = append(matches, line)
				}
			}
		}
		return matches
	}

	if len(files) == 0 {
		if stdin == "" {
			return err("grep: no input")
		}
		return ok(strings.Join(search(stdin, ""), "\n"))
	}

	var all []string
	for _, f := range files {
		n, _, e := sh.fs.resolve(sh.cwd, f)
		if e != nil {
			return err("grep: " + e.Error())
		}
		src := ""
		if len(files) > 1 {
			src = f
		}
		all = append(all, search(n.content, src)...)
	}
	if len(all) == 0 {
		return result{isErr: true} // grep returns 1 when no match (no output)
	}
	return ok(strings.Join(all, "\n"))
}

// ── wc ─────────────────────────────────────────────────────────
func cmdWC(sh *shell, args []string, stdin string) result {
	content := stdin
	name := ""
	if len(args) > 0 {
		n, _, e := sh.fs.resolve(sh.cwd, args[0])
		if e != nil {
			return err("wc: " + e.Error())
		}
		content = n.content
		name = args[0]
	}
	lines := len(strings.Split(content, "\n"))
	words := len(strings.Fields(content))
	chars := utf8.RuneCountInString(content)
	if name != "" {
		return ok(fmt.Sprintf("%7d %7d %7d %s", lines, words, chars, name))
	}
	return ok(fmt.Sprintf("%7d %7d %7d", lines, words, chars))
}

// ── find ───────────────────────────────────────────────────────
func cmdFind(sh *shell, args []string, _ string) result {
	start := sh.cwd
	nameFilter := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			nameFilter = args[i+1]
			nameFilter = strings.Trim(nameFilter, "*")
			i++
		} else if !strings.HasPrefix(args[i], "-") {
			start = args[i]
		}
	}

	root, abs, e := sh.fs.resolve(sh.cwd, start)
	if e != nil {
		return err("find: " + e.Error())
	}

	var results []string
	var walk func(n *node, p string)
	walk = func(n *node, p string) {
		match := nameFilter == "" || strings.Contains(n.name, nameFilter)
		if match {
			results = append(results, p)
		}
		if n.kind == kindDir {
			for _, c := range sh.fs.sortedEntries(n) {
				walk(c, p+"/"+c.name)
			}
		}
	}
	walk(root, abs)
	return ok(strings.Join(results, "\n"))
}

// ── man ────────────────────────────────────────────────────────
var manPages = map[string]string{
	"ls":       "LS(1)\n\nNAME\n  ls — list directory contents\n\nSYNOPSIS\n  ls [-la] [path]\n\nFLAGS\n  -l  long format (mode, size, date)\n  -a  include hidden files (dotfiles)",
	"cat":      "CAT(1)\n\nNAME\n  cat — concatenate and display files\n\nSYNOPSIS\n  cat [file...]\n\nDESCRIPTION\n  Reads files sequentially, writing them to standard output.",
	"grep":     "GREP(1)\n\nNAME\n  grep — search file for a pattern\n\nSYNOPSIS\n  grep [-i] pattern [file]\n\nFLAGS\n  -i  case-insensitive match",
	"cd":       "CD(1)\n\nNAME\n  cd — change working directory\n\nSYNOPSIS\n  cd [directory]\n\n  ~ or no argument returns to /home/anandu",
	"neofetch": "NEOFETCH(1)\n\nNAME\n  neofetch — system information tool\n\nDESCRIPTION\n  Displays system info alongside an ASCII logo.\n  This implementation shows portfolio information.",
}

func cmdMan(_ *shell, args []string, _ string) result {
	if len(args) == 0 {
		return err("man: what manual page do you want?")
	}
	page, ok2 := manPages[args[0]]
	if !ok2 {
		return err(fmt.Sprintf("man: no manual entry for %s", args[0]))
	}
	return ok(page)
}

// ── uname ──────────────────────────────────────────────────────
func cmdUname(_ *shell, args []string, _ string) result {
	if len(args) > 0 && args[0] == "-a" {
		return ok("PortfolioOS 1.0.0 portfolio #1 SMP WASM Go/js anandu.dev")
	}
	return ok("PortfolioOS")
}

// ── date ───────────────────────────────────────────────────────
func cmdDate(_ *shell, _ []string, _ string) result {
	// Note: time.Now() in WASM returns browser time
	return ok(time.Now().Format("Mon Jan  2 15:04:05 MST 2006"))
}

// ── neofetch ───────────────────────────────────────────────────
func cmdNeofetch(_ *shell, _ []string, _ string) result {
	return ok(`
    /\      anandu@portfolio
   /  \     ──────────────────────────
  / /\ \    OS:       PortfolioOS (WASM/Go)
 / ____ \   Shell:    gosh 1.0.0
/_/    \_\  Runtime:  Go → WebAssembly
            Location: Montréal, QC
            Role:     Fullstack & DevOps Engineer
            Stack:    Go, Rust, TypeScript, K8s, Docker
            Current:  Consultant @ Deloitte Canada
            Web:      anandu.dev
            Contact:  anandu_dev@icloud.com
`)
}

// ── ./exec ─────────────────────────────────────────────────────
func (sh *shell) cmdExec(cmd string, args []string) result {
	n, _, e := sh.fs.resolve(sh.cwd, cmd)
	if e != nil {
		return err(cmd + ": No such file or directory")
	}
	if n.kind == kindDir {
		return err(cmd + ": is a directory")
	}
	// Treat the file as a help text when run with --help or no args
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		return ok(n.content)
	}
	return err(fmt.Sprintf("%s: unknown argument: %s", cmd, args[0]))
}

// ── history ────────────────────────────────────────────────────
func (sh *shell) cmdHistory() result {
	if len(sh.history) == 0 {
		return ok("(empty)")
	}
	var lines []string
	for i, h := range sh.history {
		lines = append(lines, fmt.Sprintf("  %3d  %s", i+1, h))
	}
	return ok(strings.Join(lines, "\n"))
}
