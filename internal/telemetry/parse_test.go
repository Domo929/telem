package telemetry

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := os.Open("../../test/jsonStream.txt.gz")
	require.NoError(t, err)

	defer f.Close()

	r, err := gzip.NewReader(f)
	require.NoError(t, err)

	lines, err := parse(r)
	require.NoError(t, err)

	assert.Len(t, lines, 16)
}
