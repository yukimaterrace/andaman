package andaman.price

import kotlin.test.Test
import kotlin.test.assertEquals

class PriceTest {

    @Test
    fun testCurrencySymbolFileName() {
        assertEquals("USDJPY", CurrencySymbol.UsdJpy.fileName())
        assertEquals("EURJPY", CurrencySymbol.EurJpy.fileName())
        assertEquals("GBPJPY", CurrencySymbol.GbpJpy.fileName())
        assertEquals("EURUSD", CurrencySymbol.EurUsd.fileName())
        assertEquals("GBPUSD", CurrencySymbol.GbpUsd.fileName())
        assertEquals("EURGBP", CurrencySymbol.EurGbp.fileName())
    }
}
