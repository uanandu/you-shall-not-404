package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"
)

type nodeKind uint8

const (
	kindDir  nodeKind = iota
	kindFile
)

type node struct {
	name    string
	kind    nodeKind
	content string
	mtime   time.Time
	mode    string
	entries map[string]*node // nil for files
}

func dir(name string, children ...*node) *node {
	n := &node{
		name:    name,
		kind:    kindDir,
		mode:    "drwxr-xr-x",
		mtime:   buildTime,
		entries: make(map[string]*node),
	}
	for _, c := range children {
		n.entries[c.name] = c
	}
	return n
}

func file(name, content string) *node {
	return &node{
		name:    name,
		kind:    kindFile,
		mode:    "-rw-r--r--",
		mtime:   buildTime,
		content: strings.TrimLeft(content, "\n"),
	}
}

var buildTime = time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC)

type vfs struct {
	root *node
}

func (v *vfs) resolve(cwd, p string) (*node, string, error) {
	abs := v.abs(cwd, p)
	n, err := v.walk(abs)
	return n, abs, err
}

func (v *vfs) walk(abs string) (*node, error) {
	if abs == "/" {
		return v.root, nil
	}
	parts := strings.Split(strings.TrimPrefix(abs, "/"), "/")
	cur := v.root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if cur.kind != kindDir {
			return nil, fmt.Errorf("not a directory")
		}
		child, ok := cur.entries[part]
		if !ok {
			return nil, fmt.Errorf("%s: no such file or directory", abs)
		}
		cur = child
	}
	return cur, nil
}

func (v *vfs) abs(cwd, p string) string {
	if p == "" || p == "~" {
		return "/home/anandu"
	}
	if strings.HasPrefix(p, "~/") {
		p = "/home/anandu/" + p[2:]
	}
	if !strings.HasPrefix(p, "/") {
		p = cwd + "/" + p
	}
	return path.Clean(p)
}

func (v *vfs) sortedEntries(n *node) []*node {
	entries := make([]*node, 0, len(n.entries))
	for _, c := range n.entries {
		entries = append(entries, c)
	}
	sort.Slice(entries, func(i, j int) bool {
		// dirs first, then files, both alphabetical
		if entries[i].kind != entries[j].kind {
			return entries[i].kind == kindDir
		}
		return entries[i].name < entries[j].name
	})
	return entries
}
