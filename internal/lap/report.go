package lap

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

// Report gives the complete information of a driver's lap
type Report struct {
	CommonLapReport
	Sectors map[string]Sector `json:"Sectors"`
}

// isFull returns whether the struct members for a Report are filled and valid
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

// Combine takes in an initial, base, and new Report and returns an updated Report.
// The 'base' report is the report for that driver's lap. New is any subsequent data that is aggregated in before the next lap
// initial is used to transfer things like 'RacingNumber' which is set once and not referenced at any other point.
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

// FillBack handles filling data in a lap that might not be found in Combine. For example, if a driver is in the same
// position over many laps, that position will be referenced once when they gain that position, but won't be referenced
// again until it changes. This function checks each driver's lap and retroactively fills what data it can.
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
