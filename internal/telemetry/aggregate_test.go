package telemetry

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregate(t *testing.T) {
	f, err := os.Open("../../test/full.txt.gz")
	require.NoError(t, err)

	r, err := gzip.NewReader(f)
	require.NoError(t, err)

	laps, err := parse(r)
	require.NoError(t, err)

	updated, err := aggregate(laps)
	require.NoError(t, err)

	assert.Len(t, updated, 57)

	// We only look at laps 3-57, and only at Max Verstappen (33) since he had a clean race.
	// Initial laps don't have all data (such as Best Lap Time, or Last Lap Time)
	// If you are lapped, your "GapToLeader" becomes 0
	for lapNum := 2; lapNum < len(updated); lapNum++ {
		lapReport := updated[lapNum].Drivers["33"]

		assert.NotEmptyf(t, lapReport.GapToLeader, "GapToLeader|Lap ndx %d", lapNum)
		assert.NotEmptyf(t, lapReport.IntervalToPositionAhead.Value, "IntervalToPositionAhead|Lap ndx %d", lapNum)
		assert.NotEmptyf(t, lapReport.Position, "Position|Lap ndx %d", lapNum)
		assert.NotEmptyf(t, lapReport.RacingNumber, "RacingNumber|Lap ndx %d", lapNum)
		assert.NotEmptyf(t, lapReport.BestLapTime.Value, "BestLapTime|Lap ndx %d", lapNum)
		assert.NotEmptyf(t, lapReport.LastLapTime.Value, "LastLapTime|Lap ndx %d", lapNum)
	}
}
