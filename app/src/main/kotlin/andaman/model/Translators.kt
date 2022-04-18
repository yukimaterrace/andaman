package andaman.model

/**
 * ポジション情報をPosition Entityに変換します
 */
fun Position.translate(target: andaman.usecase.Position, trade: Trade) {
    positionId = target.id
    currencyPair = target.currencyPair
    amount = target.amount
    openPrice = target.openPrice.bid
    openAt = target.openAt.dbFormat()
    closePrice = target.closePrice?.bid
    closeAt = target.closeAt?.dbFormat()
    status = target.status
    profit = target.profit
    this.trade = trade
}