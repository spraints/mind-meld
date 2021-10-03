package ui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spraints/mind-meld/lmsp"
)

type targetRender struct {
	file   fileReader
	target lmsp.ProjectTarget
	blocks *blockIndex
	pos    int
}

type blockIndex struct {
	nodes    map[lmsp.ProjectBlockID]bool
	next     map[lmsp.ProjectBlockID]lmsp.ProjectBlockID
	children map[lmsp.ProjectBlockID][]lmsp.ProjectBlockID
	roots    []lmsp.ProjectBlockID
}

func (t targetRender) Init() tea.Cmd {
	return t.index
}

func (t targetRender) index() tea.Msg {
	// id -> is root?
	nodes := map[lmsp.ProjectBlockID]bool{}
	next := map[lmsp.ProjectBlockID]lmsp.ProjectBlockID{}
	children := map[lmsp.ProjectBlockID][]lmsp.ProjectBlockID{}
	for id, block := range t.target.Blocks {
		switch block := block.(type) {
		case *lmsp.ProjectBlockObject:
			if block.Next != nil {
				next[id] = *block.Next
				nodes[*block.Next] = false
			}
			if block.Parent != nil {
				children[*block.Parent] = append(children[*block.Parent], id)
			} else {
				if _, ok := nodes[id]; !ok {
					nodes[id] = true
				}
			}
		default:
			panic(fmt.Sprintf("%q %T %#v", id, block, block))
		}
	}
	var roots []lmsp.ProjectBlockID
	for id, isRoot := range nodes {
		if isRoot {
			roots = append(roots, id)
		}
	}
	sort.Slice(roots, func(a, b int) bool { return string(roots[a]) < string(roots[b]) })
	t.blocks = &blockIndex{
		nodes:    nodes,
		next:     next,
		children: children,
		roots:    roots,
	}
	return t
}

func (t targetRender) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case targetRender:
		return msg, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return t, tea.Quit
		case "ctrl+h", "left":
			return t.file, nil
		case "down":
			if t.blocks != nil && t.pos+1 < len(t.blocks.roots) {
				t.pos++
			}
			return t, nil
		case "up":
			if t.pos > 0 {
				t.pos--
			}
			return t, nil
		case "enter":
			return t.renderBlock(), nil
		}
	}
	return t, nil
}

func (t targetRender) View() string {
	if t.blocks == nil {
		return escape + loading
	}

	lines := make([]string, 0, 1+len(t.blocks.roots))
	lines = append(lines, escape)
	for i, id := range t.blocks.roots {
		if i == t.pos {
			lines = append(lines, fmt.Sprintf("> %s >\n", id))
		} else {
			lines = append(lines, fmt.Sprintf("  %s\n", id))
		}
	}
	return strings.Join(lines, "")
}

func (t targetRender) renderBlock() tea.Model {
	if t.blocks == nil || t.pos >= len(t.blocks.roots) {
		return t
	}
	rootID := t.blocks.roots[t.pos]
	visited := map[lmsp.ProjectBlockID]struct{}{}
	var res strings.Builder
	t.renderChain(&res, rootID, "> ", "  ", visited)
	return blockRender{t: t, s: res.String()}
}

func (t targetRender) renderChain(w io.Writer, id lmsp.ProjectBlockID, myIndent, followingIndent string, visited map[lmsp.ProjectBlockID]struct{}) {
	if id == "" {
		return
	}
	if _, visited := visited[id]; visited {
		fmt.Fprintf(w, "%s%s LOOP!\n", myIndent, id)
		return
	}
	visited[id] = struct{}{}

	block := t.target.Blocks[id].(*lmsp.ProjectBlockObject)
	fmt.Fprintf(w, "%s%s (%s)\n", myIndent, id, block.Opcode)

	nextID := t.blocks.next[id]

	for _, childID := range t.blocks.children[id] {
		if childID == nextID {
			continue
		}
		t.renderChain(w, childID, followingIndent+"+ ", followingIndent+"  ", visited)
	}

	t.renderChain(w, nextID, followingIndent, followingIndent, visited)
}

type blockRender struct {
	t targetRender
	s string
}

func (b blockRender) Init() tea.Cmd { return nil }
func (b blockRender) View() string  { return escape + b.s }

func (b blockRender) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "ctrl+h", "left":
			return b.t, nil
		}
	}
	return b, nil
}
