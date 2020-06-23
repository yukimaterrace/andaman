package recorder

import (
	"encoding/csv"
	"io"
	"log"
	"yukimaterrace/andaman/config"
)

// Recorder is an interface for trade recorder
type Recorder interface {
	Start()

	Record(material *Material)
}

type rowCreator interface {
	header() []string
	row(material *Material) []string
}

type routine struct {
	material chan *Material
	creator  rowCreator
	writer   io.WriteCloser
	archive  *Material
}

func newRoutine(creator rowCreator, writer io.WriteCloser) *routine {
	return &routine{
		material: make(chan *Material, config.RecorderChanCapacity),
		creator:  creator,
		writer:   writer,
	}
}

func (routine *routine) Start() {
	go routine.run()
}

func (routine *routine) Record(material *Material) {
	routine.material <- material
}

func (routine *routine) run() {
	csv := csv.NewWriter(routine.writer)

	defer routine.writer.Close()
	defer csv.Flush()

	csv.Write(routine.creator.header())

	for {
		material := <-routine.material

		row := routine.creator.row(material)
		if err := csv.Write(row); err != nil {
			log.Fatal(err)
		}

		csv.Flush()
	}
}
