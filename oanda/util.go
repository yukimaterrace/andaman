package oanda

import "fmt"

// MakeCandleSpecs is a helper method to make candle specifications
func MakeCandleSpecs(granularity string, instruments ...string) []string {
	specs := make([]string, 0)
	for _, instrument := range instruments {
		spec := fmt.Sprintf("%s:%s:M", instrument, granularity)
		specs = append(specs, spec)
	}
	return specs
}
