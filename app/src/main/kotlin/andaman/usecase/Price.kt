package andaman.usecase

import andaman.enum.BuySellType
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
 * 売買タイプからOPEN時の価格を求めます
 */
fun Price.resolveValueForOpen(buySellType: BuySellType) = when (buySellType) {
    BuySellType.BUY -> ask
    BuySellType.SELL -> bid
}

/**
 * 売買タイプからCLOSE時の価格を求めます
 */
fun Price.resolveValueForClose(buySellType: BuySellType) = when (buySellType) {
    BuySellType.BUY -> bid
    BuySellType.SELL -> ask
}

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
