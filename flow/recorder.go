package flow

type recorder interface {
	record(material recordMaterial)
	close()
}

// RecorderFactory provides factory method of recorder
type RecorderFactory interface {
	create() recorder
}

type recordMaterial interface{}

type recordRequest struct {
	material recordMaterial
}

type recordWorker struct {
	recorder
	ch chan interface{}
}

func (recordWorker *recordWorker) recordRequest(material recordMaterial) {
	recordWorker.ch <- &recordRequest{material: material}
}

func (recordWorker *recordWorker) shutdown() {
	req := newShutdownRequest()
	recordWorker.ch <- req
	<-req.done
}

func (recordWorker *recordWorker) work(exit chan<- bool) {
	request := <-recordWorker.ch

	switch req := request.(type) {
	case *recordRequest:
		recordWorker.record(req.material)
		exit <- false

	case *shutdownRequest:
		recordWorker.close()
		req.done <- true
		exit <- true
	}
}

type summarySpec interface {
	csvHeaders() []string
	csvValues() []string
}
