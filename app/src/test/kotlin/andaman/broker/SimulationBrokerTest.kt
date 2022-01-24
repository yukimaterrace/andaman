package andaman.broker

import andaman.common.Context
import andaman.price.CurrencySymbol
import andaman.price.Price
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.*
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class SimulationBrokerTest {

    private var usdJpyPrices = listOf(
        Price(CurrencySymbol.UsdJpy, BigDecimal("113.1"), LocalDateTime.of(2022, 1, 24, 12, 1)),
        Price(CurrencySymbol.UsdJpy, BigDecimal("113.2"), LocalDateTime.of(2022, 1, 24, 12, 2))
    )

    private var eurUsdPrices = listOf(
        Price(CurrencySymbol.EurUsd, BigDecimal("1.124"), LocalDateTime.of(2022, 1, 24, 12, 0)),
        Price(CurrencySymbol.EurUsd, BigDecimal("1.123"), LocalDateTime.of(2022, 1, 24, 12, 1))
    )

    private val currentPrices = usdJpyPrices.zip(eurUsdPrices).map { mapOf(it.first.symbol to it.first, it.second.symbol to it.second) }

    private val broker = SimulationBroker(Context().also { it.currentPrices = currentPrices[0] })

    @Test
    fun testNull() {
        assertNull(broker.makeMarketOrder(CurrencySymbol.EurGbp, BigDecimal("1")))
        assertNull(broker.closeOrder(UUID.randomUUID()))
    }

    @Test
    fun testOrders() {
        val position1 = broker.makeMarketOrder(CurrencySymbol.UsdJpy, BigDecimal.valueOf(10000))

        assertNotNull(position1)
        assertEquals(CurrencySymbol.UsdJpy, position1.symbol)
        assertEquals(BigDecimal.valueOf(10000), position1.amount)
        assertEquals(usdJpyPrices[0], position1.openPrice)
        assertEquals(usdJpyPrices[0].at, position1.openAt)
        assertEquals(PositionStatus.OPEN, position1.status)
        assertEquals(BigDecimal("1000.0"), position1.profit(usdJpyPrices[1]))

        val position2 = broker.makeMarketOrder(CurrencySymbol.EurUsd, BigDecimal.valueOf(100000))
        assertNotNull(position2)

        var positions = broker.positions()
        assertEquals(2, positions.size)
        assertEquals(setOf(position1.id, position2.id), positions.map { it.id }.toSet())

        positions = broker.positions(CurrencySymbol.UsdJpy)
        assertEquals(1, positions.size)
        assertEquals(position1.id, positions[0].id)

        broker.context.currentPrices = currentPrices[1]
        broker.closeOrder(position2.id)

        positions = broker.positions(CurrencySymbol.EurUsd)
        assertEquals(1, positions.size)
        assertEquals(position2.id, positions[0].id)

        val position = positions[0]
        assertEquals(CurrencySymbol.EurUsd, position.symbol)
        assertEquals(BigDecimal.valueOf(100000), position.amount)
        assertEquals(eurUsdPrices[0], position.openPrice)
        assertEquals(eurUsdPrices[0].at, position.openAt)
        assertEquals(eurUsdPrices[1], position.closePrice)
        assertEquals(eurUsdPrices[1].at, position.closeAt)
        assertEquals(PositionStatus.CLOSED, position.status)
        assertEquals(BigDecimal("-100.000"), position.profit)
    }
}
