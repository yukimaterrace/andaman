package flow

import "yukimaterrace/andaman/broker"

type simpleRecorder struct {
	writer   writer
	orderMap map[broker.AccountID]map[broker.OrderID]*completableOrder
}

func newSimpleRecorder(writer writer) *simpleRecorder {
	return &simpleRecorder{
		writer:   writer,
		orderMap: map[broker.AccountID]map[broker.OrderID]*completableOrder{},
	}
}

func (recorder *simpleRecorder) record(material recordMaterial) {
	accountCombinedOrders, ok := material.(*accountCombinedOrders)
	if !ok {
		panic("wrong type has been passed")
	}

	for accountID, combinedOrders := range accountCombinedOrders.ordersMap {
		if _, ok := recorder.orderMap[accountID]; !ok {
			recorder.orderMap[accountID] = map[broker.OrderID]*completableOrder{}
		}

		completableOrderMap := recorder.orderMap[accountID]
		for _, createdOrder := range combinedOrders.createdOrders {
			if _, ok := completableOrderMap[createdOrder.OrderID()]; ok {
				panic("duplicate order id for created order detected")
			}

			completableOrderMap[createdOrder.OrderID()] = &completableOrder{
				createdOrder: createdOrder,
				closedOrder:  nil,
			}
		}

		for _, closedOrder := range combinedOrders.closedOrders {
			completableOrder, ok := completableOrderMap[closedOrder.OrderID()]
			if !ok {
				panic("no created order for the order id exist")
			}

			completableOrder.closedOrder = closedOrder
		}
	}
}

func (recorder *simpleRecorder) write() {
	identifiedCompletableOrders := recorder.makeIdentifiedCompletableOrders(true)
	recorder.writer.write(identifiedCompletableOrders)
}

func (recorder *simpleRecorder) close() {
	identifiedCompletableOrders := recorder.makeIdentifiedCompletableOrders(false)
	recorder.writer.write(identifiedCompletableOrders)
	recorder.writer.close()
}

func (recorder *simpleRecorder) makeIdentifiedCompletableOrders(onlyCompleted bool) []*identifiedCompletableOrder {

	return nil
}

type completableOrder struct {
	createdOrder broker.CreatedOrder
	closedOrder  broker.ClosedOrder
}

type identifiedCompletableOrder struct {
	accountID broker.AccountID
	tradePair broker.TradePair
	order     *completableOrder
}

type writer interface {
	write(orders []*identifiedCompletableOrder) error
	close() error
}
