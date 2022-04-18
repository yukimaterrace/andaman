package andaman.model

import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.*

/**
 * テーブルリスト
 */
private val tables = arrayOf(Users, Trades, Positions)

/**
 * Users Table
 */
object Users: IntIdTable() {
    val accountId = integer("account_id").uniqueIndex()
    val name = name("name")
}

/**
 * Trades Table
 */
object Trades: IntIdTable() {
    val tradeId = uuid("trade_id").uniqueIndex()
    val name = name("name")

    val user = reference("user", Users)
}

/**
 * Positions Table
 */
object Positions: IntIdTable() {
    val positionId = uuid("position_id").uniqueIndex()
    val currencyPair = enumeration("currency_pair", CurrencyPair::class)
    val amount = quantity("amount")
    val openPrice = price("open_price")
    val openAt = time("open_at")
    val closePrice = price("close_price").nullable()
    val closeAt = time("close_at").nullable()
    val status = enumeration("status", PositionStatus::class)
    val profit = quantity("profit").nullable()

    val trade = reference("trade", Trades)
}

/**
 * データベースを初期化します
 */
fun Transaction.initDb(withLog: Boolean = false) {
    if (withLog) {
        addLogger(StdOutSqlLogger)
    }
    SchemaUtils.create(*tables)
}

/**
 * テーブルをドロップします
 */
fun dropTables() {
    SchemaUtils.drop(*tables)
}

private fun Table.name(name: String) = varchar(name, 50)
private fun Table.price(name: String) = decimal(name, 6, 3)
private fun Table.quantity(name: String) = decimal(name, 7, 3)
private fun Table.time(name: String) = varchar(name, 20)
