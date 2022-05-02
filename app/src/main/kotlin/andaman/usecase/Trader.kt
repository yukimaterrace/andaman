package andaman.usecase

import andaman.usecase.strategy.OrderProposal

/**
 * トレーダー
 */
interface Trader {
    /**
     * ブローカー
     */
    val broker: Broker

    /**
     * トレードを実施します
     */
    fun trade(context: Context)
}

class TraderImpl(override val broker: Broker): Trader {

    override fun trade(context: Context) {
        val proposals: List<OrderProposal> = context.strategies.map {
            it.orderProposal(context)
        }.fold(emptyList()) { acc, orderProposals ->
            acc + orderProposals
        }

        proposals.forEach {
            it.opens.forEach { openOrder ->
                broker.makeMarketOrder(context, openOrder.currencyPair, openOrder.buySellType, openOrder.amount)
            }

            it.closes.forEach { closeOrder ->
                broker.closeOrder(context, closeOrder.positionId)
            }
        }
    }
}
