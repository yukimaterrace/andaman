package andaman.model

import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import java.math.BigDecimal
import java.time.LocalDateTime

/**
 * ユーザー情報をUser Entityに変換します
 */
fun User.translate(from: andaman.usecase.User) {
    accountId = from.accountId
    name = from.name
}

/**
 * User Entityをユーザー情報に変換します
 */
fun User.toUseCase() = andaman.usecase.User(
    accountId = accountId,
    name = name
)

/**
 * トレード情報をTrade Entityに変換します
 */
fun Trade.translate(from: andaman.usecase.Trade) {
    tradeId = from.tradeId
}

/**
 * TradeEntityをトレード情報に変換します
 */
fun Trade.toUseCase() = andaman.usecase.Trade(
    tradeId = tradeId
)

/**
 * ポジション情報をPosition Entityに変換します
 */
fun Position.translate(from: andaman.usecase.Position, trade: Trade) {
    positionId = from.id
    currencyPair = from.currencyPair
    buySellType = from.buySellType
    amount = from.amount
    openPrice = from.openPrice
    openAt = from.openAt.dbFormat()
    closePrice = from.closePrice
    closeAt = from.closeAt?.dbFormat()
    status = from.status
    profit = from.profit
    this.trade = trade
}

/**
 * Position Entityをポジション情報に変換します
 */
fun Position.toUseCase() = andaman.usecase.Position(
    id = positionId,
    currencyPair = currencyPair,
    buySellType = buySellType,
    amount = amount,
    openPrice = openPrice,
    openAt = openAt.toLocalDateTime(),
    closePrice = closePrice,
    closeAt = closeAt?.toLocalDateTime(),
    status = status
)
