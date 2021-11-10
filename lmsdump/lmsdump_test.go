package lmsdump

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spraints/mind-meld/lmsp/lmspsimple"
)

func TestSample(t *testing.T) {
	expected, err := ioutil.ReadFile("testdata/project.lms.dump")
	require.NoError(t, err)

	f, err := lmspsimple.Read("testdata/project.lms")
	require.NoError(t, err)

	var buf bytes.Buffer
	assert.NoError(t, Dump(&buf, f.Project))
	assert.Equal(t, string(expected), buf.String())
}
