package lap

// Lap contains the map of driver num to Report, this is the root struct that unmarshalls all but the first line
type Lap struct {
	Drivers  map[string]Report `json:"Lines"`
	Withheld bool              `json:"Withheld"`
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
