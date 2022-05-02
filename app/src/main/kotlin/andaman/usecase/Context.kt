package andaman.usecase

import andaman.enum.CurrencyPair
import andaman.usecase.strategy.TradeStrategy
import java.util.*

/**
 * コンテクスト
 */
data class Context(
    val user: User,
    val trade: Trade,
    val strategies: List<TradeStrategy>
) {
    var currentPrices: Map<CurrencyPair, Price> = emptyMap()
    var currentPositions: List<Position> = emptyList()
}

/**
 * ユーザー情報
 */
data class User(
    val accountId: Long,
    val name: String
)

/**
 * トレード情報
 */
data class Trade(
    val tradeId: UUID
)
