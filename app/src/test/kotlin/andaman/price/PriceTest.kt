package andaman.price

import kotlin.test.Test
import kotlin.test.assertEquals

class PriceTest {

    @Test
    fun testCurrencySymbolShow() {
        assertEquals("USDJPY", CurrencySymbol.UsdJpy.show())
        assertEquals("EURJPY", CurrencySymbol.EurJpy.show())
        assertEquals("GBPJPY", CurrencySymbol.GbpJpy.show())
        assertEquals("EURUSD", CurrencySymbol.EurUsd.show())
        assertEquals("GBPUSD", CurrencySymbol.GbpUsd.show())
        assertEquals("EURGBP", CurrencySymbol.EurGbp.show())
    }
}
