package andaman.common

import andaman.broker.Position
import andaman.price.CurrencySymbol
import andaman.price.Price

class Context(val user: User) {
    var currentPrices: Map<CurrencySymbol, Price> = emptyMap()
    var currentPositions: List<Position> = emptyList()
}

data class User(
    val accountId: Long,
    val name: String
)
