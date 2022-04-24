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
    fun findUser(accountId: Long): User?

    /**
     * トレード情報を取得します
     */
    fun findTrade(tradeId: UUID): Trade?

    /**
     * ポジション情報を取得します
     */
    fun findPosition(positionId: UUID): Position?

    /**
     * ポジションリストを永続化します
     */
    fun insertPositions(positions: List<andaman.usecase.Position>, tradeId: UUID)
}

class DbRepositoryImpl: DbRepository {
    override fun findUser(accountId: Long) =
        User.find { Users.accountId eq accountId }.singleOrNull()

    override fun findTrade(tradeId: UUID) =
        Trade.find { Trades.tradeId eq tradeId }.singleOrNull()

    override fun findPosition(positionId: UUID) =
        Position.find { Positions.positionId eq positionId }.singleOrNull()

    override fun insertPositions(positions: List<andaman.usecase.Position>, tradeId: UUID) {
        val trade = findTrade(tradeId) ?: return
        positions.forEach { Position.new { translate(it, trade = trade) } }
    }
}
