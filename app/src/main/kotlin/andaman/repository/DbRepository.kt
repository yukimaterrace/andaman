package andaman.repository

import andaman.model.*
import java.util.*

/**
 * DBレポジトリ
 */
interface DbRepository {
    /**
     * ユーザー情報を取得します
     */
    fun findUser(accountId: Int): User?

    /**
     * トレード情報を取得します
     */
    fun findTrade(tradeId: UUID): Trade?

    /**
     * ポジション情報を取得します
     */
    fun findPosition(positionId: UUID): Position?
}

class DbRepositoryImpl {
    fun findUser(accountId: Int) =
        User.find { Users.accountId eq accountId }.singleOrNull()

    fun findTrade(tradeId: UUID) =
        Trade.find { Trades.tradeId eq tradeId }.singleOrNull()

    fun findPosition(positionId: UUID) =
        Position.find { Positions.positionId eq positionId }.singleOrNull()
}
