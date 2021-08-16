package lap

import "fmt"

// InitialLap contains the map of driver num to InitialLapReport, this is the root struct that unmarshalls the first line
type InitialLap struct {
	Lines    map[string]InitialLapReport `json:"Lines"`
	Withheld bool                        `json:"Withheld"`
}

// InitialLapReport gives the complete information on the initial dump of every driver's starter lap.
type InitialLapReport struct {
	CommonLapReport
	Sectors []Sector `json:"Sectors"`
}

// ToLap converts the InitialLap struct to the Lap struct
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
