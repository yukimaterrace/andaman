package andaman

import (
	"strconv"
)

type (
	oandaFloat64 float64
	oandaInt     int
	oandaBool    bool
)

func (f oandaFloat64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 10, 64)
}

func (i oandaInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (b oandaBool) String() string {
	return strconv.FormatBool(bool(b))
}
