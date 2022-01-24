package andaman.broker

import andaman.common.Context
import andaman.price.CurrencySymbol
import java.math.BigDecimal
import java.util.*

class SimulationBroker(override val context: Context): Broker {

    private var orderMap: Map<UUID, Position> = emptyMap()

    override fun makeMarketOrder(symbol: CurrencySymbol, amount: BigDecimal): Position? {
        val openPrice = context.currentPrices[symbol] ?: return null
        return Position(
            id = UUID.randomUUID(),
            symbol = symbol,
            amount = amount,
            openPrice = openPrice,
            openAt = openPrice.at
        ).also { orderMap = orderMap.plus(Pair(it.id, it)) }
    }

    override fun closeOrder(positionId: UUID): Position? {
        return orderMap[positionId]?.let { position ->
            context.currentPrices[position.symbol]?.let {
                position.closePrice = it
                position.closeAt = it.at
                position.status = PositionStatus.CLOSED
                position
            }
        }
    }

    override fun positions(symbol: CurrencySymbol?): List<Position> =
        orderMap.values.toList().let {
            symbol ?: return it
            it.filter { position -> position.symbol == symbol }
        }
}
