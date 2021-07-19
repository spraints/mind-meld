package lmsp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRoundTrip(t *testing.T) {
	f, err := os.Open("testdata/project.json")
	require.NoError(t, err)
	defer f.Close()

	raw, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	var loaded Project
	require.NoError(t, json.Unmarshal(raw, &loaded))

	rewritten, err := json.Marshal(&loaded)
	require.NoError(t, err)

	var expected map[string]interface{}
	var actual map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &expected))
	require.NoError(t, json.Unmarshal(rewritten, &actual))

	assert.Equal(t, expected, actual)
}
