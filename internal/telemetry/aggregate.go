package telemetry

func aggregate(laps []Lap) ([]Lap, error) {
	lapCount := make(map[string]int)
	for driver := range laps[0].Drivers {
		lapCount[driver] = 1
	}

	// preallocate 100 laps for speed
	aggregatedLaps := make([]Lap, 100, len(laps))
	aggregatedLaps[0] = laps[0]

	for lapNdx, lap := range laps {
		for driverNum, lapReport := range lap.Drivers {
			driverLap := lapCount[driverNum]
			if aggregatedLaps[driverLap].Drivers == nil {
				aggregatedLaps[driverLap].Drivers = make(map[string]LapReport)
			}

			updatedLap, lapDone := combine(
				aggregatedLaps[0].Drivers[driverNum],
				aggregatedLaps[driverLap].Drivers[driverNum],
				lapReport)

			aggregatedLaps[driverLap].Drivers[driverNum] = updatedLap

			if lapDone {
				updatedLap = fillBack(updatedLap, laps, lapNdx, driverNum)
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

	// fill missing data from last lap
	aggregatedLapsLastLapNdx := len(aggregatedLaps) - 1
	rawLapsLastLapNdx := len(laps) - 1
	for driverNum, lapReport := range aggregatedLaps[aggregatedLapsLastLapNdx].Drivers {
		updatedLap := fillBack(lapReport, laps, rawLapsLastLapNdx, driverNum)
		aggregatedLaps[aggregatedLapsLastLapNdx].Drivers[driverNum] = updatedLap
	}

	return aggregatedLaps, nil
}
