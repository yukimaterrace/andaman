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
    val openPrice: Price,
    val openAt: LocalDateTime
) {
    var closePrice: Price? = null
    var closeAt: LocalDateTime? = null
    var status = PositionStatus.OPEN

    val profit: BigDecimal?
        get() = closePrice?.let { calcProfit(it) }

    fun profit(price: Price) = calcProfit(price)

    private fun calcProfit(price: Price) =
        when (buySellType) {
            BuySellType.BUY -> (price.bid - openPrice.ask) * amount
            BuySellType.SELL -> (openPrice.bid - price.ask) * amount
        }
}
