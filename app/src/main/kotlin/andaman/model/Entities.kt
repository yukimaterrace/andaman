package andaman.model

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID

class User(id: EntityID<Int>): IntEntity(id) {
    companion object : IntEntityClass<User>(Users)
    var accountId by Users.accountId
    var name by Users.name
    val trades by Trade referrersOn Trades.user
}

class Trade(id: EntityID<Int>): IntEntity(id) {
    companion object : IntEntityClass<Trade>(Trades)
    var tradeId by Trades.tradeId
    var name by Trades.name
    var user by User referencedOn Trades.user
    val positions by Position referrersOn Positions.trade
}

class Position(id: EntityID<Int>): IntEntity(id) {
    companion object : IntEntityClass<Position>(Positions)
    var positionId by Positions.positionId
    var symbol by Positions.symbol
    var amount by Positions.amount
    var openPrice by Positions.openPrice
    var openAt by Positions.openAt
    var closePrice by Positions.closePrice
    var closeAt by Positions.closeAt
    var status by Positions.status
    var profit by Positions.profit
    var trade by Trade referencedOn Positions.trade
}
