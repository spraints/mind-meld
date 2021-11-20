# How to make a release

1. Make sure you've installed [goreleaser](https://goreleaser.com). (e.g. `brew install goreleaser/tap/goreleaser`)
1. Tag the release. (`git tag -a v0.1.0 -m "First release"`)
1. Set `GITHUB_TOKEN` env var. (`. ~/.github-token`)
1. Release it! (`env GITHUB_TOKEN=$GITHUB_TOKEN goreleaser release`)
