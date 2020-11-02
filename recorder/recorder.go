package recorder

import (
	"log"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

type (
	// Recorder is a struct for recorder
	Recorder struct {
		tradeRun   *model.TradeRun
		orderMap   map[broker.OrderID]*completableOrder
		openOrders []broker.OpenOrder
	}

	completableOrder struct {
		tradeConfiguration *model.TradeConfigurationDetail
		createdOrder       broker.CreatedOrder
		closedOrder        broker.ClosedOrder
	}
)

func newRecorder(tradeRun *model.TradeRun) *Recorder {
	return &Recorder{
		tradeRun: tradeRun,
		orderMap: map[broker.OrderID]*completableOrder{},
	}
}

// Record is a method to record
func (recorder *Recorder) Record(material flow.RecordMaterial) {
	recordMaterial, ok := material.(*trader.RecordMaterial)
	if !ok {
		panic(util.ErrWrongType)
	}

	if recorder.tradeRun.State == model.Pending {
		tradeRun, err := service.UpdateTradeRunForStart(recorder.tradeRun)
		if err != nil {
			panic(err)
		}
		recorder.tradeRun = tradeRun
	}

	recorder.openOrders = recordMaterial.OpenOrders

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

	if err := service.UpdateTradeRunForFinish(recorder.tradeRun); err != nil {
		log.Println(err)
	}
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
	for _, order := range completableOrders {
		if order.closedOrder == nil {
			var tradeDirection model.TradeDirection
			if order.createdOrder.IsLong() {
				tradeDirection = model.Long
			} else {
				tradeDirection = model.Short
			}

			err := service.AddCreatedOrder(
				recorder.tradeRun.TradeRunID,
				int(order.createdOrder.OrderID()),
				order.tradeConfiguration.TradeConfigurationID,
				order.createdOrder.Units(),
				tradeDirection,
				int(order.createdOrder.TimeAtOpen()),
				order.createdOrder.PriceAtOpen(),
			)

			if err != nil {
				log.Println(err)
			}
		} else {
			err := service.UpdateOrderForClose(
				recorder.tradeRun.TradeRunID,
				int(order.closedOrder.OrderID()),
				order.closedOrder.RealizedProfit(),
				int(order.closedOrder.TimeAtClose()),
				order.closedOrder.PriceAtClose(),
			)

			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, openOrder := range recorder.openOrders {
		err := service.UpdateOrderForProfit(
			recorder.tradeRun.TradeRunID,
			int(openOrder.OrderID()),
			openOrder.UnrealizedProfit(),
		)

		if err != nil {
			log.Println(err)
		}
	}
}
