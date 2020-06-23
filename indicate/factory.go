package indicate

// CreateSMA25 is a factory method for SMA25 indicator
func CreateSMA25() Indicator {
	return newRoutine(newSmaCalculator(SMA25))
}

// CreateSMA75 is a factory method for SMA75 indicator
func CreateSMA75() Indicator {
	return newRoutine(newSmaCalculator(SMA75))
}

// CreateSMA150 is a factory method for SMA150 indicator
func CreateSMA150() Indicator {
	return newRoutine(newSmaCalculator(SMA150))
}

// CreateSMA600 is a factory method for SMA600 indicator
func CreateSMA600() Indicator {
	return newRoutine(newSmaCalculator(SMA600))
}
