package andaman.dao

import andaman.broker.PositionStatus
import andaman.model.*
import andaman.price.CurrencySymbol
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.transactions.transaction
import java.math.BigDecimal
import java.nio.file.Paths
import java.util.*
import kotlin.random.Random
import kotlin.test.*

class DaoTest {

    @BeforeTest
    fun setup() {
        val path = Paths.get("resource/test/test.db").toAbsolutePath().toString()
        Database.connect("jdbc:sqlite:${path}", "org.sqlite.JDBC")
        transaction { initDb() }
    }

    @AfterTest
    fun shutdown() {
        transaction { dropTables() }
    }

    @Test
    fun test() {
        transaction {
            val userA = createUser("A")
            val tradeA = createTrade("A", userA)
            val position1 = createPosition(tradeA)
            val position2 = createPosition(tradeA)

            val user = findUser(userA.accountId)
            assertNotNull(user)
            assertEquals(userA.name, user.name)
            assertEquals(1, user.trades.count())
            assertEquals(tradeA.name, user.trades.first().name)

            val trade = findTrade(tradeA.tradeId)
            assertNotNull(trade)
            assertEquals(userA.accountId, trade.user.accountId)
            assertEquals(2, trade.positions.count())

            val positions = trade.positions.toList().filter { it.positionId == position1.positionId }
            assertEquals(1, positions.size)

            val position = positions.first()
            assertEquals(CurrencySymbol.UsdJpy, position.symbol)
            assertEquals(BigDecimal("10000.000"), position.amount)
            assertEquals(BigDecimal("112.500"), position.openPrice)

            position2.closePrice = BigDecimal("113.61234")
            position2.closeAt = "2022-01-27-09:05"
            position2.profit = BigDecimal("135.242")

            val position0 = findPosition(position2.positionId)
            assertNotNull(position0)
            assertEquals(BigDecimal("113.612"), position0.closePrice)
            assertEquals("2022-01-27-09:05", position0.closeAt)
            assertEquals(BigDecimal("135.242"), position0.profit)
            assertEquals(tradeA.tradeId, position0.trade.tradeId)
        }
    }

    private fun createUser(name: String) =
        User.new {
            accountId = Random.nextInt(99999999)
            this.name = name
        }

    private fun createTrade(name: String, user: User) =
        Trade.new {
            tradeId = UUID.randomUUID()
            this.name = name
            this.user = user
        }

    private fun createPosition(trade: Trade) =
        Position.new {
            positionId = UUID.randomUUID()
            symbol = CurrencySymbol.UsdJpy
            amount = BigDecimal("10000.0")
            openPrice = BigDecimal("112.5")
            openAt = "2022-01-26-12:24"
            status = PositionStatus.OPEN
            this.trade = trade
        }
}
