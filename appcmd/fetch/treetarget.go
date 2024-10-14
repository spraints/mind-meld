package fetch

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type treeTarget struct {
	repo  *git.Repository
	blobs map[string]plumbing.Hash
}

func newTreeTarget(repo *git.Repository) *treeTarget {
	return &treeTarget{
		repo:  repo,
		blobs: map[string]plumbing.Hash{},
	}
}

func (tt *treeTarget) Add(name string, data []byte) error {
	oid, err := createBlob(tt.repo, data)
	if err != nil {
		return err
	}
	tt.blobs[name] = oid
	return nil
}

func (tt *treeTarget) Finish() (plumbing.Hash, error) {
	return createTree(tt.repo, tt.blobs)
}

func createBlob(g *git.Repository, data []byte) (plumbing.Hash, error) {
	obj := g.Storer.NewEncodedObject()
	obj.SetType(plumbing.BlobObject)

	w, err := obj.Writer()
	if err != nil {
		return plumbing.ZeroHash, err
	}

	n, err := w.Write(data)
	if err != nil {
		return plumbing.ZeroHash, err
	}
	if n != len(data) {
		return plumbing.ZeroHash, fmt.Errorf("incomplete write")
	}

	if err := w.Close(); err != nil {
		return plumbing.ZeroHash, err
	}

	return g.Storer.SetEncodedObject(obj)
}

func createTree(g *git.Repository, blobs map[string]plumbing.Hash) (plumbing.Hash, error) {
	var tb treeBuilder
	for name, oid := range blobs {
		tb.Add(name, oid)
	}
	return tb.Build(g)
}
