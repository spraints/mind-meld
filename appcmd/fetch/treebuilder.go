package fetch

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const gitPathSeparator = "/"

type treeBuilder struct {
	subdirs map[string]*treeBuilder
	blobs   map[string]plumbing.Hash
}

func (tb *treeBuilder) Add(path string, oid plumbing.Hash) {
	parts := strings.SplitN(path, gitPathSeparator, 2)
	switch len(parts) {
	case 1:
		if tb.blobs == nil {
			tb.blobs = map[string]plumbing.Hash{
				path: oid,
			}
		} else {
			tb.blobs[path] = oid
		}
	case 2:
		subdir := parts[0]
		rest := parts[1]
		if tb.subdirs == nil {
			tb.subdirs = map[string]*treeBuilder{
				subdir: {},
			}
		} else if tb.subdirs[subdir] == nil {
			tb.subdirs[subdir] = &treeBuilder{}
		}
		tb.subdirs[subdir].Add(rest, oid)
	default:
		panic("illegal path " + path)
	}
}

func (tb *treeBuilder) Build(g *git.Repository) (plumbing.Hash, error) {
	var entries []object.TreeEntry
	if tb.blobs != nil {
		for name, oid := range tb.blobs {
			entries = append(entries, object.TreeEntry{
				Name: name,
				Mode: 0100644,
				Hash: oid,
			})
		}
	}
	if tb.subdirs != nil {
		for name, subdir := range tb.subdirs {
			subtree, err := subdir.Build(g)
			if err != nil {
				return plumbing.ZeroHash, fmt.Errorf("%s: %w", name, err)
			}
			entries = append(entries, object.TreeEntry{
				Name: name,
				Mode: 040000,
				Hash: subtree,
			})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })

	var t object.Tree
	t.Entries = entries

	obj := g.Storer.NewEncodedObject()
	if err := t.Encode(obj); err != nil {
		return plumbing.ZeroHash, err
	}

	return g.Storer.SetEncodedObject(obj)
}
