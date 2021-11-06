# How to use mind-meld with Git

## View local diffs with mind-meld

In your repository, add this to `.gitattributes` and check it in.

    *.lms diff=mind-meld

On your computer, set up mind-meld's git diff tool.

    $ git config --global diff.mind-meld.textconv 'mind-meld git-diff'

## Pre-commit hook: add text version of programs

If you want to have a text copy of the program alongside the programs, add a pre-commit hook to generate them.

    mind-meld pre-commit

Backfill the text files using filter-branch:

    git filter-branch --index-filter 'mind-meld pre-commit --cached'
