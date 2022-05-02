package andaman.usecase.strategy

import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.usecase.Context
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.*

/**
 * トレードストラテジー
 */
interface TradeStrategy {
    /** 対象 */
    val subject: TradeStrategySubject

    /**
     * オーダーを提唱します
     */
    fun orderProposal(context: Context): OrderProposal
}

/**
 * トレードストラテジー対象
 */
data class TradeStrategySubject(
    val currencyPair: CurrencyPair,
    val buySellTypes: Set<BuySellType>,
    val tradeTime: TradeTime
)

/**
 * トレード時間
 */
data class TradeTime(
    val startTime: LocalDateTime,
    val endTime: LocalDateTime
)

/**
 * オープンオーダー提唱
 */
data class OpenOrderProposal(
    val currencyPair: CurrencyPair,
    val buySellType: BuySellType,
    val amount: BigDecimal
)

/**
 * クローズオーダー提唱
 */
data class CloseOrderProposal(
    val positionId: UUID,
    val amount: BigDecimal? = null
)

/**
 * オーダー提唱
 */
data class OrderProposal(
    val opens: List<OpenOrderProposal>,
    val closes: List<CloseOrderProposal>
)
