package telemetry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := os.Open("../../test/jsonStream.txt")
	require.NoError(t, err)

	defer f.Close()

	lines, err := parse(f)
	require.NoError(t, err)

	assert.Len(t, lines, 16)
}
