package record

import (
	"fmt"
	"os"
	"path"
	"time"
	"yukimaterrace/andaman/config"
	"yukimaterrace/andaman/market"
)

// CreateSimpleFileRecorder is a factory method for simple file recorder
func CreateSimpleFileRecorder(instrument market.Instrument) Recorder {
	fileName := fmt.Sprintf("%s-%d", instrument, time.Now().Unix())
	path := path.Join(config.RecordFilePath, fileName)

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	return newRoutine(&simpleRowCreator{}, file)
}
