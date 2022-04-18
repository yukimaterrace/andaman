package andaman.usecase

import andaman.enum.CurrencyPair
import kotlin.test.Test
import kotlin.test.assertEquals

class PriceTest {

    @Test
    fun testCurrencySymbolFileName() {
        assertEquals("USDJPY", CurrencyPair.UsdJpy.fileName())
        assertEquals("EURJPY", CurrencyPair.EurJpy.fileName())
        assertEquals("GBPJPY", CurrencyPair.GbpJpy.fileName())
        assertEquals("EURUSD", CurrencyPair.EurUsd.fileName())
        assertEquals("GBPUSD", CurrencyPair.GbpUsd.fileName())
        assertEquals("EURGBP", CurrencyPair.EurGbp.fileName())
    }
}
