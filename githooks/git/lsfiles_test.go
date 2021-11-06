package git

import (
	"bufio"
	"bytes"
	"testing"
)

func TestParseLsFiles(t *testing.T) {
	r := bytes.NewBuffer([]byte("100644 28b7a60587b10d9d5a0f15e94596b29473bfe2fb 0\tMakefile\n" +
		"100644 a6755a65e8514deeed4c0c89da51dfbad9b9c8f6 0\tgithooks/git/gitdir.go\n"))
	expectedEntries := []IndexEntry{
		{Mode: 0100644, OID: "28b7a60587b10d9d5a0f15e94596b29473bfe2fb", Index: 0, Path: "Makefile"},
		{Mode: 0100644, OID: "a6755a65e8514deeed4c0c89da51dfbad9b9c8f6", Index: 0, Path: "githooks/git/gitdir.go"},
	}

	actualEntries, err := parseLsFiles(bufio.NewReader(r))
	if err != nil {
		t.Error(err)
	}
	if len(expectedEntries) != len(actualEntries) {
		t.Errorf("expected %d entries but got %d", len(expectedEntries), len(actualEntries))
	}
	for i := 0; i < len(expectedEntries) || i < len(actualEntries); i++ {
		if i < len(expectedEntries) {
			expected := expectedEntries[i]
			if i < len(actualEntries) {
				actual := actualEntries[i]
				if expected != actual {
					t.Errorf("entries[%d]: expected %#v but got %#v", i, expected, actual)
				}
			} else {
				t.Errorf("entries[%d]: expected %#v but got nothing", i, expected)
			}
		} else {
			t.Errorf("entries[%d]: expected nothing but got %#v", i, actualEntries[i])
		}
	}
}
