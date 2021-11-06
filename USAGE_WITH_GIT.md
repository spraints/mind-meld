# How to use mind-meld with Git

## Pre-commit hook: add text version of programs

Add this to `.git/hooks/pre-commit':

    mind-meld pre-commit

## Backfill text versions of programs

Backfill using filter-branch:

    git filter-branch --index-filter 'mind-meld pre-commit --cached'
