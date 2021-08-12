package telemetry

import "fmt"

// Line contains the map of driver num to LapReport, this is the root struct that unmarshalls all but the first line
type Line struct {
	Lines    map[string]LapReport `json:"Lines"`
	Withheld bool                 `json:"Withheld"`
}

// InitialLine contains the map of driver num to InitialLapReport, this is the root struct that unmarshalls the first line
type InitialLine struct {
	Lines    map[string]InitialLapReport `json:"Lines"`
	Withheld bool                        `json:"Withheld"`
}

// CommonLapReport contains struct fields that are in both InitialLapReport and LapReport
type CommonLapReport struct {
	GapToLeader             string            `json:"GapToLeader"`
	IntervalToPositionAhead Interval          `json:"IntervalToPositionAhead"`
	Line                    int               `json:"Line"`
	Position                string            `json:"Position"`
	ShowPosition            bool              `json:"ShowPosition"`
	RacingNumber            string            `json:"RacingNumber"`
	Retired                 bool              `json:"Retired"`
	InPit                   bool              `json:"InPit"`
	PitOut                  bool              `json:"PitOut"`
	Stopped                 bool              `json:"Stopped"`
	Status                  int               `json:"Status"`
	Sectors                 map[string]Sector `json:"Sectors"`
	Speeds                  Speeds            `json:"Speeds"`
	BestLapTime             LapTime           `json:"BestLapTime"`
	LastLapTime             LapTime           `json:"LastLapTime"`
}

// InitialLapReport gives the complete information on the initial dump of every driver's starter lap.
type InitialLapReport struct {
	CommonLapReport
	Sectors []Sector `json:"Sectors"`
}

// LapReport gives the complete information of a driver's lap
type LapReport struct {
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
	Status          int    `json:"Status"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

// Speed denotes a driver's speed at various points on the track
type Speed struct {
	Value           string `json:"Value"`
	Status          int    `json:"Status"`
	OverallFastest  bool   `json:"OverallFastest"`
	PersonalFastest bool   `json:"PersonalFastest"`
}

// Speeds holds multiple Speed structs for the different meansuring locations on a track
type Speeds struct {
	I1 Speed `json:"I1"`
	I2 Speed `json:"I2"`
	FL Speed `json:"FL"`
	ST Speed `json:"ST"`
}

// LapTime gives the information about a drivers lap time
type LapTime struct {
	Value           string `json:"Value"`
	Status          int    `json:"Status,omitempty"`
	OverallFastest  bool   `json:"OverallFastest,omitempty"`
	PersonalFastest bool   `json:"PersonalFastest,omitempty"`
}

func (i InitialLine) ConvertToLine() *Line {
	line := new(Line)
	line.Withheld = i.Withheld
	line.Lines = make(map[string]LapReport)

	for driver, initialReport := range i.Lines {
		report := LapReport{
			CommonLapReport: initialReport.CommonLapReport,
			Sectors:         make(map[string]Sector),
		}
		for ndx, sector := range initialReport.Sectors {
			report.Sectors[fmt.Sprint(ndx)] = sector
		}
		line.Lines[driver] = report
	}

	return line
}
