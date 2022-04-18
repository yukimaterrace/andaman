package andaman.usecase

import andaman.enum.CurrencyPair
import java.math.BigDecimal
import java.time.LocalDateTime

/**
 * 価格情報
 */
data class Price(
    val currencyPair: CurrencyPair,
    val bid: BigDecimal,
    val ask: BigDecimal,
    val at: LocalDateTime
)

/**
 * 通貨ペアと価格のマップ
 */
typealias CurrencyPairMap = Map<CurrencyPair, Price>

/**
 * PriceSupplier インターフェース
 */
typealias PriceSupplier = Iterable<CurrencyPairMap>

/**
 * 通貨ペアのファイル名文字列を生成します
 */
fun CurrencyPair.fileName(): String = this.name.uppercase()
