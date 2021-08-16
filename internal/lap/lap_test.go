package lap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLapReport_Combine(t *testing.T) {
	var (
		initial = Report{
			CommonLapReport: CommonLapReport{
				Position:     "1",
				RacingNumber: "77",
			},
		}
		empty = Report{}
		full  = Report{
			CommonLapReport: CommonLapReport{
				GapToLeader: "+1.1",
				IntervalToPositionAhead: Interval{
					Value:    "+1.5",
					Catching: true,
				},
				Position:     "2",
				RacingNumber: "44",
				Retired:      true,
				InPit:        true,
				PitOut:       true,
				Stopped:      true,
				Speeds: Speeds{
					FL: Speed{
						Value:           "400",
						OverallFastest:  true,
						PersonalFastest: true,
					},
					ST: Speed{
						Value:           "500",
						OverallFastest:  true,
						PersonalFastest: true,
					},
				},
				BestLapTime: Time{
					Value:           "1:23.00",
					OverallFastest:  true,
					PersonalFastest: true,
				},
				LastLapTime: Time{
					Value:           "1:25.00",
					OverallFastest:  true,
					PersonalFastest: true,
				},
			},
			Sectors: map[string]Sector{
				"1": {
					Stopped:         true,
					Value:           "24.0",
					OverallFastest:  true,
					PersonalFastest: true,
				},
				"2": {
					Stopped:         true,
					Value:           "25.0",
					OverallFastest:  true,
					PersonalFastest: true,
				},
				"3": {
					Stopped:         true,
					Value:           "26.0",
					OverallFastest:  true,
					PersonalFastest: true,
				},
			},
		}
	)

	testTable := []struct {
		name            string
		initial         Report
		base            Report
		new             Report
		expectedUpdated Report
		expectedDone    bool
	}{
		{
			name:            "empty fields",
			initial:         initial,
			base:            empty,
			new:             full,
			expectedUpdated: full,
			expectedDone:    true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			actualUpdated, actualDone := combine(test.initial, test.base, test.new)
			assert.Equal(t, test.expectedUpdated, actualUpdated)
			assert.Equal(t, test.expectedDone, actualDone)
		})
	}
}
