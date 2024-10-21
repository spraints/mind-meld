# Mind Meld - inspect LEGO Mindstorms Scratch programs

`mind-meld` exists to help you see what's changed in your Mindstorms programs.

It has two modes of operation:

1. It can fetch Python programs from the Spike or Mindstorms app and save them
   to a directory or Git branch.

2. It can show you what's changed between two block programs.

## Python

### Fetch python programs into a directory

For example, to see what changes you're making, you might do this:

```
# Create a repository for your Python programs.
$ git init my-python-programs
$ cd my-python-programs

# Fetch your current Python programs and commit them.
$ mind-meld spike fetch --dir .
$ git add .
$ git commit -m "Initial commit"

# Make some changes to your Python programs in the Spike app.

# Fetch the changes and diff them.
$ mind-meld spike fetch --dir .
$ git diff

# Commit them if you like at this point.
```

### Fetch python programs into a Git branch

You might also just want to build a Git branch that contains the changes to
your Python programs.

```
# Create a repository for your Python programs.
$ git init my-python-programs
$ cd my-python-programs

# Fetch your current Python programs and commit them.
$ mind-meld spike fetch --git refs/lego/scratch -m "Initial commit"

# Make some changes to your Python programs in the Spike app.

# Fetch the changes into a new commit on your tracking branch.
$ mind-meld spike fetch --git refs/lego/scratch -m "Update from Spike"

# See the changes.
$ git show refs/lego/scratch

# Check them out in your working directory.
$ git merge refs/lego/scratch

# Make some more changes in the Spike app.

# Fetch the changes into a new commit on your tracking branch and merge them
# again into your main branch.
$ mind-meld spike fetch --git refs/lego/scratch -m "Update from Spike"
$ git merge refs/lego/scratch

# See the latest change.
$ git show

# Make some more changes in the Spike app.

# Now, just see what's different.
$ mind-meld spike fetch --dir .
$ git diff

# If you want to reset your working directory so that you can 'fetch --git' and
# 'merge' again easily, do this:
$ git checkout .
$ git clean -fd
```

## Blocks

### View diffs with mind-meld

In your repository, add this to `.gitattributes` and check it in.

    *.lms diff=mind-meld
    *.lmsp diff=mind-meld

On your computer, install mind-meld and set up mind-meld's git diff tool.

    $ go install github.com/spraints/mind-meld@latest # or download a release
    $ git config --global diff.mind-meld.textconv 'mind-meld dump'

Now, git log will look like this:

![git log example](docs/git-log.png)
