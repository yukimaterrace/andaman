package andaman.broker

import andaman.common.Context
import andaman.price.CurrencySymbol
import andaman.price.Price
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.*

interface Broker {
    val context: Context

    fun makeMarketOrder(symbol: CurrencySymbol, amount: BigDecimal): Position?
    fun closeOrder(positionId: UUID): Position?
    fun positions(symbol: CurrencySymbol? = null): List<Position>
}

class Position(
    val id: UUID,
    val symbol: CurrencySymbol,
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

    private fun calcProfit(price: Price) = (price.value - openPrice.value) * amount
}

enum class PositionStatus {
    OPEN,
    CLOSED
}
