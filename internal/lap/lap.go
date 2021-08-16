package lap

import (
	"fmt"
)

// Lap contains the map of driver num to Report, this is the root struct that unmarshalls all but the first line
type Lap struct {
	Drivers  map[string]Report `json:"Lines"`
	Withheld bool              `json:"Withheld"`
}

// InitialLap contains the map of driver num to InitialLapReport, this is the root struct that unmarshalls the first line
type InitialLap struct {
	Lines    map[string]InitialLapReport `json:"Lines"`
	Withheld bool                        `json:"Withheld"`
}

// CommonLapReport contains struct fields that are in both InitialLapReport and Report
type CommonLapReport struct {
	GapToLeader             string   `json:"GapToLeader"`
	IntervalToPositionAhead Interval `json:"IntervalToPositionAhead"`
	Position                string   `json:"Position"`
	RacingNumber            string   `json:"RacingNumber"`
	Retired                 bool     `json:"Retired"`
	InPit                   bool     `json:"InPit"`
	PitOut                  bool     `json:"PitOut"`
	Stopped                 bool     `json:"Stopped"`
	Speeds                  Speeds   `json:"Speeds"`
	BestLapTime             Time     `json:"BestLapTime"`
	LastLapTime             Time     `json:"LastLapTime"`
}

// InitialLapReport gives the complete information on the initial dump of every driver's starter lap.
type InitialLapReport struct {
	CommonLapReport
	Sectors []Sector `json:"Sectors"`
}

// Report gives the complete information of a driver's lap
type Report struct {
	CommonLapReport
	Sectors map[string]Sector `json:"Sectors"`
}

// Interval gives the time to the driver ahead,
// and a boolean about whether that time is increasing (false) or decreasing (true)
type Interval struct {
	Value    string `json:"Value"`
	Catching bool   `json:"Catching"`
}

// Sector gives information about a driver's sector split times
type Sector struct {
	Stopped         bool   `json:"Stopped"`
	Value           string `json:"Value"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

// Speed denotes a driver's speed at various points on the track
type Speed struct {
	Value           string `json:"Value"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

// Speeds holds multiple Speed structs for the different measuring locations on a track
type Speeds struct {
	FL Speed `json:"FL"`
	ST Speed `json:"ST"`
}

// Time gives the information about a drivers lap time
type Time struct {
	Value           string `json:"Value"`
	OverallFastest  bool   `json:"OverallFastest,omitempty"`
	PersonalFastest bool   `json:"PersonalFastest,omitempty"`
}

func (i InitialLap) ToLap() *Lap {
	line := new(Lap)
	line.Withheld = i.Withheld
	line.Drivers = make(map[string]Report)

	for driver, initialReport := range i.Lines {
		report := Report{
			CommonLapReport: initialReport.CommonLapReport,
			Sectors:         make(map[string]Sector),
		}
		for ndx, sector := range initialReport.Sectors {
			report.Sectors[fmt.Sprint(ndx)] = sector
		}
		line.Drivers[driver] = report
	}

	return line
}

func (lr Report) isFull() bool {
	return lr.RacingNumber != "" &&
		lr.Position != "" &&
		lr.GapToLeader != "" &&
		lr.IntervalToPositionAhead.Value != "" &&
		lr.Speeds.ST.Value != "" &&
		lr.Speeds.FL.Value != "" &&
		lr.BestLapTime != Time{} &&
		lr.LastLapTime != Time{}
}

func Combine(initial, base, new Report) (Report, bool) {
	// make a local copy of base to avoid overwriting
	updated := base

	updated.RacingNumber = initial.RacingNumber
	// These just get updated to the latest value
	if new.GapToLeader != "" {
		updated.GapToLeader = new.GapToLeader
	}

	if new.IntervalToPositionAhead.Value != "" {
		updated.IntervalToPositionAhead = new.IntervalToPositionAhead
	}

	if new.Position != "" {
		updated.Position = new.Position
	}

	if !updated.Retired && new.Retired {
		updated.Retired = new.Retired
	}
	updated.InPit = updated.InPit || new.InPit
	updated.PitOut = updated.PitOut || new.PitOut
	updated.Stopped = updated.Stopped || new.Stopped

	if new.BestLapTime.Value != "" {
		updated.BestLapTime = new.BestLapTime
	}

	if new.LastLapTime.Value != "" {
		updated.LastLapTime = new.LastLapTime
	}

	// These get merged in a more intelligent manner
	if updated.Speeds.FL.Value == "" {
		updated.Speeds.FL = new.Speeds.FL
	}
	if updated.Speeds.ST.Value == "" {
		updated.Speeds.ST = new.Speeds.ST
	}

	if updated.Sectors == nil {
		updated.Sectors = make(map[string]Sector)
	}
	for sectorNum, sector := range new.Sectors {
		if _, ok := updated.Sectors[sectorNum]; !ok {
			if sector.Value != "" {
				updated.Sectors[sectorNum] = sector
			}
		}
	}

	// return the updated lapReport, as well as if the lap is 'done' (all three sectors exist and have time)
	return updated, len(updated.Sectors) == 3
}

func FillBack(lap Report, rawLaps []Lap, lapNdx int, driverNum string) Report {
	for i := lapNdx; i >= 0; i-- {
		if lap.isFull() {
			return lap
		}

		rawLap := rawLaps[i].Drivers[driverNum]

		if lap.Position == "" && rawLap.Position != "" {
			lap.Position = rawLap.Position
		}

		if lap.BestLapTime.Value == "" && rawLap.BestLapTime.Value != "" {
			lap.BestLapTime = rawLap.BestLapTime
		}

		if lap.LastLapTime.Value == "" && rawLap.LastLapTime.Value != "" {
			lap.LastLapTime = rawLap.LastLapTime
		}
	}

	return lap
}
