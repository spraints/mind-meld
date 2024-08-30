package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spraints/mind-meld/appcmd"
	"github.com/spraints/mind-meld/lmsp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Run(app appcmd.App) error {
	g, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	projects, err := listProjects(app)

	blobs := map[string]plumbing.Hash{}
	for _, project := range projects {
		data, err := readPythonProject(project)
		if err != nil {
			return fmt.Errorf("%s: %w", project.Name, err)
		}

		if data != nil {
			oid, err := createBlob(g, data)
			if err != nil {
				return fmt.Errorf("%s: %w", project.Name, err)
			}
			blobs[project.Name+".py"] = oid
		}
	}

	tree, err := createTree(g, blobs)
	if err != nil {
		return err
	}

	return createCommit(g, plumbing.ReferenceName(app.RemoteName()), tree,
		"Update copy of "+app.FullName()+" python programs")
}

type project struct {
	Name string
	Path string
}

func listProjects(app appcmd.App) ([]project, error) {
	for _, d := range app.ProjectDirs() {
		found, err := walkProjectDir(d, "", nil)
		if err == nil {
			return found, nil
		}
	}
	return nil, fmt.Errorf("no project dir found (checked %v)", app.ProjectDirs())
}

func walkProjectDir(dirname string, prefix string, result []project) ([]project, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			p, err := walkProjectDir(filepath.Join(dirname, e.Name()), prefix+e.Name()+"/", result)
			if err != nil {
				return nil, err
			}
			result = p
			continue
		}

		if e.Type().IsRegular() {
			result = append(result, project{
				Name: prefix + e.Name(),
				Path: filepath.Join(dirname, e.Name()),
			})
		}
	}
	return result, nil
}

func readPythonProject(proj project) ([]byte, error) {
	f, err := os.Open(proj.Path)
	if err != nil {
		return nil, err
	}

	l, err := lmsp.ReadFile(f)
	if err != nil {
		return nil, err
	}

	man, err := l.Manifest()
	if err != nil {
		return nil, err
	}

	if man.Type != "python" {
		fmt.Printf("%s: skip %s program\n", proj.Name, man.Type)
		return nil, nil
	}

	program, err := l.Python()
	if err != nil {
		return nil, err
	}
	return []byte(program), nil
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

type treeID = plumbing.Hash

func createTree(g *git.Repository, blobs map[string]plumbing.Hash) (treeID, error) {
	var tb treeBuilder
	for name, oid := range blobs {
		tb.Add(name, oid)
	}
	return tb.Build(g)
}

func createCommit(g *git.Repository, refName plumbing.ReferenceName, tree treeID, commitMsg string) error {
	// Check the ref.
	// If the tree is the same, there's nothing to do.
	// If the ref is there, use its OID as the parent commit.
	var parentHashes []plumbing.Hash
	if ref, err := g.Reference(refName, false); err == nil {
		c, err := g.CommitObject(ref.Hash())
		if err != nil {
			return err
		}
		if c.TreeHash == tree {
			return nil
		}
		parentHashes = append(parentHashes, ref.Hash())
	}

	// Get an author and committer to use for the commit.
	var o git.CommitOptions
	if err := o.Validate(g); err != nil {
		return err
	}

	// Build the commit.
	var c object.Commit
	c.Author = *o.Author
	c.Committer = *o.Committer
	c.Message = commitMsg
	c.TreeHash = tree
	c.ParentHashes = parentHashes

	// Save the commit.
	obj := g.Storer.NewEncodedObject()
	if err := c.Encode(obj); err != nil {
		return err
	}

	commitID, err := g.Storer.SetEncodedObject(obj)
	if err != nil {
		return err
	}

	// Update the reference.
	return g.Storer.SetReference(plumbing.NewHashReference(refName, commitID))
}
