package telemetry

import "github.com/Domo929/telem.git/internal/lap"

func aggregate(laps []lap.Lap) ([]lap.Lap, error) {
	lapCount := make(map[string]int)
	for driver := range laps[0].Drivers {
		lapCount[driver] = 1
	}

	// preallocate 100 laps for speed
	aggregatedLaps := make([]lap.Lap, 100, len(laps))
	aggregatedLaps[0] = laps[0]

	for lapNdx, rawLap := range laps {
		for driverNum, lapReport := range rawLap.Drivers {
			driverLap := lapCount[driverNum]
			if aggregatedLaps[driverLap].Drivers == nil {
				aggregatedLaps[driverLap].Drivers = make(map[string]lap.Report)
			}

			updatedLap, lapDone := lap.Combine(
				aggregatedLaps[0].Drivers[driverNum],
				aggregatedLaps[driverLap].Drivers[driverNum],
				lapReport)

			aggregatedLaps[driverLap].Drivers[driverNum] = updatedLap

			if lapDone {
				updatedLap = lap.FillBack(updatedLap, laps, lapNdx, driverNum)
				aggregatedLaps[driverLap].Drivers[driverNum] = updatedLap

				lapCount[driverNum] += 1

			}
		}
	}

	// remove the trailing entries that are not valid laps
	for i, lap := range aggregatedLaps {
		if len(lap.Drivers) == 0 {
			aggregatedLaps = aggregatedLaps[:i]
			break
		}
	}

	// fill missing data from last rawLap
	aggregatedLapsLastLapNdx := len(aggregatedLaps) - 1
	rawLapsLastLapNdx := len(laps) - 1
	for driverNum, lapReport := range aggregatedLaps[aggregatedLapsLastLapNdx].Drivers {
		updatedLap := lap.FillBack(lapReport, laps, rawLapsLastLapNdx, driverNum)
		aggregatedLaps[aggregatedLapsLastLapNdx].Drivers[driverNum] = updatedLap
	}

	return aggregatedLaps, nil
}
