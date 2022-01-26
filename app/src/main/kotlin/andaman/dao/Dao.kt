package andaman.dao

import andaman.model.*
import java.util.*

fun findUser(accountId: Int) =
    User.find { Users.accountId eq accountId }.singleOrNull()

fun findTrade(tradeId: UUID) =
    Trade.find { Trades.tradeId eq tradeId }.singleOrNull()

fun findPosition(positionId: UUID) =
    Position.find { Positions.positionId eq positionId }.singleOrNull()
