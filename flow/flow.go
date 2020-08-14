package flow

type worker interface {
	work(exit chan<- bool)
}

type runner struct {
	worker
	exit chan bool
}

func newRunner(worker worker) *runner {
	return &runner{
		worker: worker,
		exit:   make(chan bool, 1),
	}
}

func (runner *runner) run(done chan<- bool) {
	go func() {
		for {
			runner.work(runner.exit)

			if <-runner.exit {
				break
			}
		}
		done <- true
	}()
}

// Flow is a struct for flow
type Flow struct {
	priceWorker  *priceWorker
	tradeWorker  *tradeWorker
	recordWorker *recordWorker

	priceWorkerDone  chan bool
	tradeWorkerDone  chan bool
	recordWorkerDone chan bool
}

// Start is a method to start the flow
func (flow *Flow) Start() {
	newRunner(flow.priceWorker).run(flow.priceWorkerDone)
	newRunner(flow.tradeWorker).run(flow.tradeWorkerDone)
	newRunner(flow.recordWorker).run(flow.recordWorkerDone)
}

// WaitCompletion waits until shutdown of the flow
func (flow *Flow) WaitCompletion() {
	<-flow.priceWorkerDone
	<-flow.tradeWorkerDone
	<-flow.recordWorkerDone
}

// Shutdown makes shutdown of the flow
func (flow *Flow) Shutdown() {
	flow.priceWorker.shutdown()
}

// ChangeTradeMode change trade mode of the flow
func (flow *Flow) ChangeTradeMode(mode TradeMode) {
	flow.tradeWorker.changeTradeModeRequest(mode)
}

type shutdownRequest struct {
	done chan bool
}

func newShutdownRequest() *shutdownRequest {
	return &shutdownRequest{done: make(chan bool, 1)}
}
