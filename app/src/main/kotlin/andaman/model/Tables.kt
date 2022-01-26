package andaman.model

import andaman.broker.PositionStatus
import andaman.price.CurrencySymbol
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.*

private val tables = arrayOf(Users, Trades, Positions)

object Users: IntIdTable() {
    val accountId = integer("account_id").uniqueIndex()
    val name = name("name")
}

object Trades: IntIdTable() {
    val tradeId = uuid("trade_id").uniqueIndex()
    val name = name("name")

    val user = reference("user", Users)
}

object Positions: IntIdTable() {
    val positionId = uuid("position_id").uniqueIndex()
    val symbol = enumeration("symbol", CurrencySymbol::class)
    val amount = quantity("amount")
    val openPrice = price("open_price")
    val openAt = time("open_at")
    var closePrice = price("close_price").nullable()
    var closeAt = time("close_at").nullable()
    var status = enumeration("status", PositionStatus::class)
    val profit = quantity("profit").nullable()

    val trade = reference("trade", Trades)
}

fun Transaction.initDb(withLog: Boolean = false) {
    if (withLog) {
        addLogger(StdOutSqlLogger)
    }
    SchemaUtils.create(*tables)
}

fun dropTables() {
    SchemaUtils.drop(*tables)
}

private fun Table.name(name: String) = varchar(name, 50)
private fun Table.price(name: String) = decimal(name, 6, 3)
private fun Table.quantity(name: String) = decimal(name, 7, 3)
private fun Table.time(name: String) = varchar(name, 20)
