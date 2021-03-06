package flow

import "time"

type (
	// Recorder is an interface for recorder
	Recorder interface {
		Record(material RecordMaterial)
		Write()
		Close()
	}

	// RecorderFactory provides factory method of recorder
	RecorderFactory interface {
		Create() Recorder
	}

	// RecordMaterial is a interface for record material
	RecordMaterial interface{}

	recordRequest struct {
		material RecordMaterial
	}
)

type recordWorker struct {
	Recorder
	ch     chan interface{}
	ticker *time.Ticker
}

func (recordWorker *recordWorker) recordRequest(material RecordMaterial) {
	recordWorker.ch <- &recordRequest{material: material}
}

func (recordWorker *recordWorker) shutdown() {
	req := newShutdownRequest()
	recordWorker.ch <- req
	<-req.done
}

func (recordWorker *recordWorker) work(exit chan<- bool) {
	select {
	case request := <-recordWorker.ch:

		switch req := request.(type) {
		case *recordRequest:
			recordWorker.Record(req.material)

		case *shutdownRequest:
			recordWorker.ticker.Stop()
			recordWorker.Close()
			req.done <- true
			exit <- true
		}

	case <-recordWorker.ticker.C:
		recordWorker.Write()
	}
}
