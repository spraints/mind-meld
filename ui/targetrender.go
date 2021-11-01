package ui

import (
	"fmt"
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
	objs     map[lmsp.ProjectBlockID]*lmsp.ProjectBlockObject
}

func (t targetRender) Init() tea.Cmd {
	return t.index
}

func (t targetRender) index() tea.Msg {
	// id -> is root?
	nodes := map[lmsp.ProjectBlockID]bool{}
	next := map[lmsp.ProjectBlockID]lmsp.ProjectBlockID{}
	children := map[lmsp.ProjectBlockID][]lmsp.ProjectBlockID{}
	objs := map[lmsp.ProjectBlockID]*lmsp.ProjectBlockObject{}
	for id, block := range t.target.Blocks {
		switch block := block.(type) {
		case *lmsp.ProjectBlockObject:
			objs[id] = block
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
		objs:     objs,
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
			lines = append(lines, fmt.Sprintf("> %s %s\n", id, describe(t.blocks.objs[id])))
		} else {
			lines = append(lines, fmt.Sprintf("  %s %s\n", id, describe(t.blocks.objs[id])))
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
	lines := t.renderChain(nil, rootID, "> ", "  ", visited)
	return blockRender{t: t, lines: lines}
}

func (t targetRender) renderChain(res []string, id lmsp.ProjectBlockID, myIndent, followingIndent string, visited map[lmsp.ProjectBlockID]struct{}) []string {
	if id == "" {
		return res
	}
	if _, visited := visited[id]; visited {
		return append(res, fmt.Sprintf("%s%s LOOP!\n", myIndent, id))
	}
	visited[id] = struct{}{}

	res = append(res, fmt.Sprintf("%s%s %s\n", myIndent, id, describe(t.blocks.objs[id])))

	nextID := t.blocks.next[id]

	for _, childID := range t.blocks.children[id] {
		if childID == nextID {
			continue
		}
		res = t.renderChain(res, childID, followingIndent+"+ ", followingIndent+"  ", visited)
	}

	return t.renderChain(res, nextID, followingIndent, followingIndent, visited)
}

// todo - move this to the lmsp package?
func describe(block *lmsp.ProjectBlockObject) string {
	if block == nil {
		return "NIL??"
	}
	return describeOpcode(block.Opcode) + describeInputs(block.Inputs) + describeFields(block.Fields)
}

var human = map[lmsp.ProjectOpcode]string{
	"flipperevents_whenProgramStarts": "(when program starts)",
	"flippersensors_resetYaw":         "Reset yaw",
	//"flippermotor_motorGoDirectionToPosition": "
}

func describeOpcode(opcode lmsp.ProjectOpcode) string {
	if s, ok := human[opcode]; ok {
		return s
	}
	return string(opcode)
}

func describeInputs(inputs lmsp.TODO) string {
	if inputs != nil {
		return ""
	}
	return fmt.Sprintf(" %#v", inputs)
}

func describeFields(fields map[lmsp.ProjectFieldName]lmsp.ProjectField) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf(" %v", fields)
}

// todo - get this from the window size!
const blockRenderLines = 10

type blockRender struct {
	t     targetRender
	lines []string
	pos   int
}

func (b blockRender) Init() tea.Cmd { return nil }

func (b blockRender) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "ctrl+h", "left":
			return b.t, nil
		case "down":
			if b.pos+blockRenderLines < len(b.lines) {
				b.pos++
			}
		case "up":
			if b.pos > 0 {
				b.pos--
			}
		}
	}
	return b, nil
}

func (b blockRender) View() string {
	var head, tail = "\n", "\n"
	lines := b.lines[b.pos:]
	if len(lines) > blockRenderLines {
		lines = lines[:blockRenderLines]
		tail = "vvvvv\n"
	}
	if b.pos > 0 {
		head = "^^^^^\n"
	}
	return escape + head + strings.Join(lines, "") + tail
}
