package andaman.usecase

import andaman.enum.CurrencyPair
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull

class ProductsHolderTest {

    @Test
    fun testProductsHolder() {
        val sourcePath = "resource/test/test_products.yaml"
        val holder = ProductsHolder(sourcePath)
        holder.load()

        val usdJpy = holder.product(CurrencyPair.UsdJpy)
        assertNotNull(usdJpy)
        assertEquals(CurrencyPair.UsdJpy, usdJpy.currencyPair)
        assertEquals("0.01".toBigDecimal(), usdJpy.pipsUnit)
        assertEquals("0.2".toBigDecimal(), usdJpy.simulationSpreadPips)

        val eurUsd = holder.product(CurrencyPair.EurUsd)
        assertNotNull(eurUsd)
        assertEquals(CurrencyPair.EurUsd, eurUsd.currencyPair)
        assertEquals("0.0001".toBigDecimal(), eurUsd.pipsUnit)
        assertEquals("0.4".toBigDecimal(), eurUsd.simulationSpreadPips)
    }
}
