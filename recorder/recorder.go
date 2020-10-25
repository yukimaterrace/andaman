package recorder

import (
	"sort"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

// Recorder is a struct for recorder
type Recorder struct {
	writer   writer
	orderMap map[trader.PartitionID]map[broker.OrderID]*completableOrder
}

func newRecorder(writer writer) *Recorder {
	return &Recorder{
		writer:   writer,
		orderMap: map[trader.PartitionID]map[broker.OrderID]*completableOrder{},
	}
}

// Record is a method to record
func (recorder *Recorder) Record(material flow.RecordMaterial) {
	partitionCombinedOrders, ok := material.(trader.PartitionCombinedOrders)
	if !ok {
		panic(util.ErrWrongType)
	}

	for partitionID, combinedOrders := range partitionCombinedOrders {
		if _, ok := recorder.orderMap[partitionID]; !ok {
			recorder.orderMap[partitionID] = map[broker.OrderID]*completableOrder{}
		}

		completableOrderMap := recorder.orderMap[partitionID]
		for _, createdOrder := range combinedOrders.CreatedOrders {
			if _, ok := completableOrderMap[createdOrder.OrderID()]; ok {
				panic("duplicate order id for created order detected")
			}

			completableOrderMap[createdOrder.OrderID()] = &completableOrder{
				createdOrder: createdOrder,
				closedOrder:  nil,
			}
		}

		for _, closedOrder := range combinedOrders.ClosedOrders {
			completableOrder, ok := completableOrderMap[closedOrder.OrderID()]
			if !ok {
				panic("no created order for the order id exist")
			}

			completableOrder.closedOrder = closedOrder
		}
	}
}

// Write is a method to write
func (recorder *Recorder) Write() {
	identifiedCompletableOrders := recorder.flush(true)
	recorder.writer.write(identifiedCompletableOrders)
}

// Close is a method to close
func (recorder *Recorder) Close() {
	identifiedCompletableOrders := recorder.flush(false)

	recorder.writer.write(identifiedCompletableOrders)
	recorder.writer.close()
}

func (recorder *Recorder) flush(onlyCompleted bool) identifiedCompletableOrders {
	orders := identifiedCompletableOrders{}

	for partitionID, orderMap := range recorder.orderMap {
		var closedOrderIDs []broker.OrderID

		for orderID, order := range orderMap {
			if onlyCompleted && order.closedOrder == nil {
				continue
			}

			orders = append(orders, &identifiedCompletableOrder{
				partitionID: partitionID,
				tradePair:   order.createdOrder.TradePair(),
				order:       order,
			})

			if order.closedOrder != nil {
				closedOrderIDs = append(closedOrderIDs, orderID)
			}
		}

		for _, orderID := range closedOrderIDs {
			delete(orderMap, orderID)
		}
	}

	sort.Sort(orders)
	return orders
}

type (
	completableOrder struct {
		createdOrder broker.CreatedOrder
		closedOrder  broker.ClosedOrder
	}

	identifiedCompletableOrder struct {
		partitionID trader.PartitionID
		tradePair   broker.TradePair
		order       *completableOrder
	}
)

func (identifiedCompletableOrder *identifiedCompletableOrder) csvHeaders() []string {
	return []string{
		"orderID",
		"tradePair",
		"timeAtOpen",
		"priceAtOpen",
		"units",
		"isLong",
		"timeAtClose",
		"priceAtClose",
		"realizedPL",
	}
}

func (identifiedCompletableOrder *identifiedCompletableOrder) csvValues() []string {
	created := identifiedCompletableOrder.order.createdOrder
	closed := identifiedCompletableOrder.order.closedOrder

	csv := []string{
		strconv.FormatInt(int64(created.OrderID()), 10),
		string(created.TradePair()),
		strconv.FormatInt(created.TimeAtOpen(), 10),
		strconv.FormatFloat(created.PriceAtOpen(), 'f', 6, 64),
		strconv.FormatFloat(created.Units(), 'f', 8, 64),
		strconv.FormatBool(created.IsLong()),
	}

	if closed == nil {
		return append(csv, "not closed", "not closed", "0")
	}

	return append(csv,
		strconv.FormatInt(closed.TimeAtClose(), 10),
		strconv.FormatFloat(closed.PriceAtClose(), 'f', 6, 64),
		strconv.FormatFloat(closed.RealizedProfit(), 'f', 6, 64),
	)
}

type identifiedCompletableOrders []*identifiedCompletableOrder

func (orders identifiedCompletableOrders) Len() int {
	return len(orders)
}

func (orders identifiedCompletableOrders) Less(i, j int) bool {
	return orders[i].order.createdOrder.OrderID() < orders[j].order.createdOrder.OrderID()
}

func (orders identifiedCompletableOrders) Swap(i, j int) {
	order := orders[i]
	orders[i] = orders[j]
	orders[j] = order
}
