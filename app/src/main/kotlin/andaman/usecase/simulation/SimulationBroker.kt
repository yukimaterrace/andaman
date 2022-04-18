package andaman.usecase.simulation

import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import andaman.usecase.Broker
import andaman.usecase.Context
import andaman.usecase.Position
import java.math.BigDecimal
import java.util.*

/**
 * シミュレーション用のブローカー
 */
class SimulationBroker: Broker {

    private var orderMap: Map<UUID, Position> = emptyMap()

    /**
     * マーケットオーダーを発行します
     */
    override fun makeMarketOrder(context: Context, currencyPair: CurrencyPair, buySellType: BuySellType, amount: BigDecimal): Position? {
        val openPrice = context.currentPrices[currencyPair] ?: return null
        return Position(
            id = UUID.randomUUID(),
            currencyPair = currencyPair,
            buySellType = buySellType,
            amount = amount,
            openPrice = openPrice,
            openAt = openPrice.at
        ).also { orderMap = orderMap.plus(Pair(it.id, it)) }
    }

    /**
     * ポジションをクローズするオーダーを発行します
     */
    override fun closeOrder(context: Context, positionId: UUID): Position? {
        return orderMap[positionId]?.let { position ->
            context.currentPrices[position.currencyPair]?.let {
                position.closePrice = it
                position.closeAt = it.at
                position.status = PositionStatus.CLOSED
                position
            }
        }
    }

    /**
     * ポジション一覧を取得します
     */
    override fun positions(symbol: CurrencyPair?): List<Position> =
        orderMap.values.toList().let {
            symbol ?: return it
            it.filter { position -> position.currencyPair == symbol }
        }
}
