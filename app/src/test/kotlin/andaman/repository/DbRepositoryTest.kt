package andaman.repository

import andaman.TestDbAware
import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import andaman.model.Position
import andaman.model.Trade
import andaman.model.User
import andaman.testPosition
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*
import kotlin.random.Random
import kotlin.test.*

class DbRepositoryTest: TestDbAware {

    private val repository = DbRepositoryImpl()

    @BeforeTest
    fun setup() {
        setupTestDb()
    }

    @AfterTest
    fun shutdown() {
        shutdownTestDb()
    }

    @Test
    fun testFind() {
        transaction {
            val userA = createUser("A")
            val tradeA = createTrade("A", userA)
            val position1 = createPosition(tradeA)
            val position2 = createPosition(tradeA)

            val user = repository.findUser(userA.accountId)
            assertNotNull(user)
            assertEquals(userA.name, user.name)
            assertEquals(1, user.trades.count())
            assertEquals(tradeA.name, user.trades.first().name)

            val trade = repository.findTrade(tradeA.tradeId)
            assertNotNull(trade)
            assertEquals(userA.accountId, trade.user.accountId)
            assertEquals(2, trade.positions.count())

            val positions = trade.positions.toList().filter { it.positionId == position1.positionId }
            assertEquals(1, positions.size)

            val position = positions.first()
            assertEquals(CurrencyPair.UsdJpy, position.currencyPair)
            assertEquals(BuySellType.BUY, position.buySellType)
            assertEquals("10000.000".toBigDecimal(), position.amount)
            assertEquals("112.500".toBigDecimal(), position.openPrice)

            position2.closePrice = "113.61234".toBigDecimal()
            position2.closeAt = "2022-01-27-09:05"
            position2.profit = "135.242".toBigDecimal()

            val position0 = repository.findPosition(position2.positionId)
            assertNotNull(position0)
            assertEquals("113.612".toBigDecimal(), position0.closePrice)
            assertEquals("2022-01-27-09:05", position0.closeAt)
            assertEquals("135.242".toBigDecimal(), position0.profit)
            assertEquals(tradeA.tradeId, position0.trade.tradeId)
        }
    }

    @Test
    fun testInsertPositions() {
        transaction {
            val user = createUser("B")
            val trade = createTrade("B", user)
            val positions = listOf(testPosition())
            val positionId = positions.first().id

            repository.insertPositions(positions, trade.tradeId)
            val position = repository.findPosition(positionId)

            assertNotNull(position)
            assertEquals(positionId, position.positionId)
            assertEquals(trade.tradeId, position.trade.tradeId)
        }
    }

    private fun createUser(name: String) =
        User.new {
            accountId = Random.nextLong(99999999)
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
            currencyPair = CurrencyPair.UsdJpy
            buySellType = BuySellType.BUY
            amount = "10000.0".toBigDecimal()
            openPrice = "112.5".toBigDecimal()
            openAt = "2022-01-26-12:24"
            status = PositionStatus.OPEN
            this.trade = trade
        }
}
