package andaman.price

import java.math.BigDecimal
import java.nio.file.Paths
import java.time.LocalDateTime
import kotlin.io.path.Path
import kotlin.test.Test
import kotlin.test.assertEquals

class SimulationPriceSupplierTest {

    @Test
    fun testCurrencySymbolPriceFilePaths() {
        val start = LocalDateTime.of(2022, 9, 15, 12, 11)
        val end = LocalDateTime.of(2023, 1, 10, 12, 11)
        val actual = CurrencySymbol.UsdJpy.priceFilePaths(start, end, "path")
        val expected = listOf(
                Pair(CurrencySymbol.UsdJpy, Path("path", "USDJPY_2022_09.csv")),
                Pair(CurrencySymbol.UsdJpy, Path("path", "USDJPY_2022_10.csv")),
                Pair(CurrencySymbol.UsdJpy, Path("path", "USDJPY_2022_11.csv")),
                Pair(CurrencySymbol.UsdJpy, Path("path", "USDJPY_2022_12.csv")),
                Pair(CurrencySymbol.UsdJpy, Path("path", "USDJPY_2023_01.csv")),
        )
        assertEquals(expected, actual)
    }

    @Test
    fun testStringPrice() {
        val s = "2022.01.20,18:06,121.0,122.0,123.0,124.0,99"
        val actual = s.currencySymbolPrice(CurrencySymbol.UsdJpy)
        val expected = Pair(
            CurrencySymbol.UsdJpy,
            Price(CurrencySymbol.UsdJpy, BigDecimal("124.0"), LocalDateTime.of(2022, 1, 20, 18, 6))
        )
        assertEquals(expected, actual)
    }

    @Test
    fun testSimulationPriceSupplier() {
        val sp = SimulationPriceSupplier(
            start = LocalDateTime.of(2022, 2, 28, 23, 58),
            end = LocalDateTime.of(2022, 3, 1, 0, 3),
            symbols = setOf(CurrencySymbol.UsdJpy, CurrencySymbol.EurUsd),
            filePath = Paths.get("resource/test").toAbsolutePath().toString()
        )
        val actual = sp.toList()
        val expected = listOf(
            mapOf(
                CurrencySymbol.EurUsd to Price(CurrencySymbol.EurUsd, BigDecimal("0.613"), LocalDateTime.of(2022, 2, 28, 23, 58))
            ),
            mapOf(
                CurrencySymbol.UsdJpy to Price(CurrencySymbol.UsdJpy, BigDecimal("123.234"), LocalDateTime.of(2022, 2, 28, 23, 59)),
                CurrencySymbol.EurUsd to Price(CurrencySymbol.EurUsd, BigDecimal("0.713"), LocalDateTime.of(2022, 2, 28, 23, 59))
            ),
            mapOf(
                CurrencySymbol.UsdJpy to Price(CurrencySymbol.UsdJpy, BigDecimal("113.234"), LocalDateTime.of(2022, 3, 1, 0, 0)),
            ),
            mapOf(
                CurrencySymbol.UsdJpy to Price(CurrencySymbol.UsdJpy, BigDecimal("133.234"), LocalDateTime.of(2022, 3, 1, 0, 1)),
                CurrencySymbol.EurUsd to Price(CurrencySymbol.EurUsd, BigDecimal("0.913"), LocalDateTime.of(2022, 3, 1, 0, 1))
            ),
            mapOf(
                CurrencySymbol.EurUsd to Price(CurrencySymbol.EurUsd, BigDecimal("0.510"), LocalDateTime.of(2022, 3, 1, 0, 2))
            ),
            mapOf(
                CurrencySymbol.UsdJpy to Price(CurrencySymbol.UsdJpy, BigDecimal("154.0"), LocalDateTime.of(2022, 3, 1, 0, 3))
            )
        )
        assertEquals(expected, actual)
    }
}
