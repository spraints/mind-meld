package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexItem(t *testing.T) {
	assert.Equal(t, "100644,0123456789abcdef0123456789abcdef01234567,example/path.txt",
		indexItem(0100644, "0123456789abcdef0123456789abcdef01234567", "example/path.txt"))
}
