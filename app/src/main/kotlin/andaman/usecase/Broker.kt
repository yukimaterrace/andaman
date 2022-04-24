package andaman.usecase

import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.*

/**
 * ブローカー
 */
interface Broker {
    fun makeMarketOrder(context: Context, currencyPair: CurrencyPair, buySellType: BuySellType, amount: BigDecimal): Position?
    fun closeOrder(context: Context, positionId: UUID): Position?
    fun positions(symbol: CurrencyPair? = null): List<Position>
}

/**
 * ポジション
 */
class Position(
    val id: UUID,
    val currencyPair: CurrencyPair,
    val buySellType: BuySellType,
    val amount: BigDecimal,
    val openPrice: BigDecimal,
    val openAt: LocalDateTime,
    var closePrice: BigDecimal? = null,
    var closeAt: LocalDateTime? = null,
    var status: PositionStatus = PositionStatus.OPEN
) {
    val profit: BigDecimal?
        get() = closePrice?.profit()

    fun profit(price: Price) = price.resolveValueForProfit().profit()

    private fun BigDecimal.profit(): BigDecimal {
        val value = (this - openPrice) * amount
        return when (buySellType) {
            BuySellType.BUY -> value
            BuySellType.SELL -> -value
        }
    }

    private fun Price.resolveValueForProfit(): BigDecimal =
        when (buySellType) {
            BuySellType.BUY -> bid
            BuySellType.SELL -> ask
        }
}
