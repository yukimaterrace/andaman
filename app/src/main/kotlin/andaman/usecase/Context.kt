package andaman.usecase

import andaman.enum.CurrencyPair
import org.kodein.di.DI
import org.kodein.di.DIAware
import java.util.*

/**
 * コンテクスト
 */
class Context(
    val user: User,
    val trade: Trade,
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
