package andaman.price

import java.math.BigDecimal
import java.time.LocalDateTime

data class Price
(
    val symbol: CurrencySymbol,
    val value: BigDecimal,
    val at: LocalDateTime
)

typealias PriceMap = Map<CurrencySymbol, Price>
typealias PriceSupplier = Iterable<PriceMap>

enum class CurrencySymbol {
    UsdJpy,
    EurJpy,
    GbpJpy,
    EurUsd,
    GbpUsd,
    EurGbp
}

fun CurrencySymbol.fileName(): String = this.name.uppercase()
