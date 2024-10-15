package fetch

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitTarget struct {
	Ref           string
	CommitMessage string
}

func (t GitTarget) refName() plumbing.ReferenceName {
	if strings.HasPrefix(string(t.Ref), "refs/") {
		return plumbing.ReferenceName(string(t.Ref))
	}
	return plumbing.ReferenceName("refs/heads/" + string(t.Ref))
}

func (t GitTarget) commitMessage() string {
	if t.CommitMessage != "" {
		return t.CommitMessage
	}
	return "Update copy of mindstorms programs"
}

func (t GitTarget) Open() (TargetInstance, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}
	return &gitTargetInstance{t, repo, NewTreeBuilder(repo)}, nil
}

func (t GitTarget) PathSeparator() string {
	return GitPathSeparator
}

type gitTargetInstance struct {
	dest GitTarget
	repo *git.Repository
	tt   *TreeBuilder
}

func (g *gitTargetInstance) Add(name string, data []byte) error {
	return g.tt.Add(name, data)
}

func (g *gitTargetInstance) Finish() (string, error) {
	tree, err := g.tt.Finish()
	if err != nil {
		return "", err
	}

	targetRef := g.dest.refName()
	commitID, err := createCommit(g.repo, targetRef, tree, g.dest.commitMessage())
	if err != nil {
		return "", err
	}

	if commitID.IsZero() {
		return fmt.Sprintf("%s: no changes found", targetRef), nil
	}

	return fmt.Sprintf("%s: created commit %v", targetRef, commitID), nil
}

func createCommit(g *git.Repository, refName plumbing.ReferenceName, tree plumbing.Hash, commitMsg string) (plumbing.Hash, error) {
	// Check the ref.
	// If the tree is the same, there's nothing to do.
	// If the ref is there, use its OID as the parent commit.
	var parentHashes []plumbing.Hash
	if ref, err := g.Reference(refName, false); err == nil {
		c, err := g.CommitObject(ref.Hash())
		if err != nil {
			return plumbing.ZeroHash, err
		}
		if c.TreeHash == tree {
			return plumbing.ZeroHash, nil
		}
		parentHashes = append(parentHashes, ref.Hash())
	}

	// Get an author and committer to use for the commit.
	var o git.CommitOptions
	if err := o.Validate(g); err != nil {
		return plumbing.ZeroHash, err
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
		return plumbing.ZeroHash, err
	}

	commitID, err := g.Storer.SetEncodedObject(obj)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	if os.Getenv("FOO") == "FOO" {
		w, _ := g.Worktree()
		w.Commit("", nil)
	}

	// Update the reference.
	return commitID, g.Storer.SetReference(plumbing.NewHashReference(refName, commitID))
}
