package recorder

import (
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

type (
	// Recorder is a struct for recorder
	Recorder struct {
		orderMap       map[broker.OrderID]*completableOrder
		priceExtractor broker.PriceExtractor
	}

	completableOrder struct {
		tradeConfiguration *model.TradeConfigurationDetail
		createdOrder       broker.CreatedOrder
		closedOrder        broker.ClosedOrder
	}
)

// NewRecorder is a constructor for recorder
func NewRecorder() *Recorder {
	return &Recorder{
		orderMap: map[broker.OrderID]*completableOrder{},
	}
}

// Record is a method to record
func (recorder *Recorder) Record(material flow.RecordMaterial) {
	recordMaterial, ok := material.(trader.RecordMaterial)
	if !ok {
		panic(util.ErrWrongType)
	}

	recorder.priceExtractor, ok = recordMaterial.TradeMaterial.(broker.PriceExtractor)
	if !ok {
		panic(util.ErrWrongType)
	}

	for _, partitionCombinedOrder := range recordMaterial.PartitionCombinedOrders {
		for _, createdOrder := range partitionCombinedOrder.CreatedOrders {
			if _, ok := recorder.orderMap[createdOrder.OrderID()]; ok {
				panic("duplicate order id for created order detected")
			}

			recorder.orderMap[createdOrder.OrderID()] = &completableOrder{
				tradeConfiguration: partitionCombinedOrder.TradeConfiguration,
				createdOrder:       createdOrder,
				closedOrder:        nil,
			}
		}

		for _, closedOrder := range partitionCombinedOrder.ClosedOrders {
			completableOrder, ok := recorder.orderMap[closedOrder.OrderID()]
			if !ok {
				panic("no created order for the order id exist")
			}

			completableOrder.closedOrder = closedOrder
		}
	}
}

// Write is a method to write
func (recorder *Recorder) Write() {
	completableOrders := recorder.flatOrders()
	recorder.flush(completableOrders)
}

// Close is a method to close
func (recorder *Recorder) Close() {
	recorder.Write()
}

func (recorder *Recorder) flatOrders() []*completableOrder {
	var orders []*completableOrder
	var closedOrderIDs []broker.OrderID

	for _, order := range recorder.orderMap {
		orders = append(orders, order)

		if order.closedOrder != nil {
			closedOrderIDs = append(closedOrderIDs, order.createdOrder.OrderID())
		}
	}

	for _, orderID := range closedOrderIDs {
		delete(recorder.orderMap, orderID)
	}

	return orders
}

func (recorder *Recorder) flush(completableOrders []*completableOrder) {

}
