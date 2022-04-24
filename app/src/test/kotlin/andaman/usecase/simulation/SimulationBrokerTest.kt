package andaman.usecase.simulation

import andaman.enum.BuySellType
import andaman.enum.CurrencyPair
import andaman.enum.PositionStatus
import andaman.testContext
import andaman.usecase.Price
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.*
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class SimulationBrokerTest {

    private var usdJpyPrices = listOf(
        Price(CurrencyPair.UsdJpy, "113.1".toBigDecimal(), "113.2".toBigDecimal(), LocalDateTime.of(2022, 1, 24, 12, 1)),
        Price(CurrencyPair.UsdJpy, "113.6".toBigDecimal(), "113.7".toBigDecimal(), LocalDateTime.of(2022, 1, 24, 12, 2))
    )

    private var eurUsdPrices = listOf(
        Price(CurrencyPair.EurUsd, "1.124".toBigDecimal(), "1.125".toBigDecimal(), LocalDateTime.of(2022, 1, 24, 12, 0)),
        Price(CurrencyPair.EurUsd, "1.129".toBigDecimal(), "1.130".toBigDecimal(), LocalDateTime.of(2022, 1, 24, 12, 1))
    )

    private val currentPrices = usdJpyPrices.zip(eurUsdPrices).map { mapOf(it.first.currencyPair to it.first, it.second.currencyPair to it.second) }

    private val broker = SimulationBroker()
    private val context = testContext().also { it.currentPrices = currentPrices[0] }

    @Test
    fun testNull() {
        assertNull(broker.makeMarketOrder(context, CurrencyPair.EurGbp, BuySellType.BUY, BigDecimal("1")))
        assertNull(broker.closeOrder(context, UUID.randomUUID()))
    }

    @Test
    fun testOrders() {
        val position1 = broker.makeMarketOrder(context, CurrencyPair.UsdJpy, BuySellType.BUY, BigDecimal.valueOf(10000))

        assertNotNull(position1)
        assertEquals(CurrencyPair.UsdJpy, position1.currencyPair)
        assertEquals(BuySellType.BUY, position1.buySellType)
        assertEquals(BigDecimal.valueOf(10000), position1.amount)
        assertEquals(usdJpyPrices[0].ask, position1.openPrice)
        assertEquals(usdJpyPrices[0].at, position1.openAt)
        assertEquals(PositionStatus.OPEN, position1.status)
        assertEquals("4000.0".toBigDecimal(), position1.profit(usdJpyPrices[1]))

        val position2 = broker.makeMarketOrder(context, CurrencyPair.EurUsd, BuySellType.SELL, BigDecimal.valueOf(100000))
        assertNotNull(position2)

        var positions = broker.positions()
        assertEquals(2, positions.size)
        assertEquals(setOf(position1.id, position2.id), positions.map { it.id }.toSet())

        positions = broker.positions(CurrencyPair.UsdJpy)
        assertEquals(1, positions.size)
        assertEquals(position1.id, positions[0].id)

        context.currentPrices = currentPrices[1]
        broker.closeOrder(context, position2.id)

        positions = broker.positions(CurrencyPair.EurUsd)
        assertEquals(1, positions.size)
        assertEquals(position2.id, positions[0].id)

        val position = positions[0]
        assertEquals(CurrencyPair.EurUsd, position.currencyPair)
        assertEquals(BuySellType.SELL, position.buySellType)
        assertEquals(BigDecimal.valueOf(100000), position.amount)
        assertEquals(eurUsdPrices[0].bid, position.openPrice)
        assertEquals(eurUsdPrices[0].at, position.openAt)
        assertEquals(eurUsdPrices[1].ask, position.closePrice)
        assertEquals(eurUsdPrices[1].at, position.closeAt)
        assertEquals(PositionStatus.CLOSED, position.status)
        assertEquals("-600.000".toBigDecimal(), position.profit)
    }
}
