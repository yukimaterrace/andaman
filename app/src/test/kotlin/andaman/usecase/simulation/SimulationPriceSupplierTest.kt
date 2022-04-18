package andaman.usecase.simulation

import andaman.enum.CurrencyPair
import andaman.testDI
import andaman.usecase.Price
import andaman.usecase.Product
import java.nio.file.Paths
import java.time.LocalDateTime
import kotlin.io.path.Path
import kotlin.test.Test
import kotlin.test.assertEquals

class SimulationPriceSupplierTest {

    @Test
    fun testProductPriceFilePaths() {
        val start = LocalDateTime.of(2022, 9, 15, 12, 11)
        val end = LocalDateTime.of(2023, 1, 10, 12, 11)
        val currencyPair = CurrencyPair.UsdJpy
        val actual = currencyPair.priceFilePaths(start, end, "path")
        val expected = listOf(
                currencyPair to Path("path", "USDJPY_2022_09.csv"),
                currencyPair to Path("path", "USDJPY_2022_10.csv"),
                currencyPair to Path("path", "USDJPY_2022_11.csv"),
                currencyPair to Path("path", "USDJPY_2022_12.csv"),
                currencyPair to Path("path", "USDJPY_2023_01.csv"),
        )
        assertEquals(expected, actual)
    }

    @Test
    fun testStringPrice() {
        val s = "2022.01.20,18:06,121.0,122.0,123.0,124.0,99"
        val currencyPair = CurrencyPair.UsdJpy
        val actual = s.currencyPairPrice(Product(currencyPair, "0.01".toBigDecimal(), "0.2".toBigDecimal()))
        val expected = currencyPair to
            Price(currencyPair, "124.0".toBigDecimal(), "124.002".toBigDecimal(), LocalDateTime.of(2022, 1, 20, 18, 6))
        assertEquals(expected, actual)
    }

    @Test
    fun testSimulationPriceSupplier() {
        val usdJpy = CurrencyPair.UsdJpy
        val eurUsd = CurrencyPair.EurUsd
        val sp = SimulationPriceSupplier(
            start = LocalDateTime.of(2022, 2, 28, 23, 58),
            end = LocalDateTime.of(2022, 3, 1, 0, 3),
            currencyPairs = setOf(usdJpy, eurUsd),
            filePath = Paths.get("resource/test").toAbsolutePath().toString(),
            di = testDI()
        )
        val actual = sp.toList()
        val expected = listOf(
            mapOf(
                eurUsd to Price(eurUsd, "0.613".toBigDecimal(), "0.61304".toBigDecimal(), LocalDateTime.of(2022, 2, 28, 23, 58))
            ),
            mapOf(
                usdJpy to Price(usdJpy, "123.234".toBigDecimal(), "123.236".toBigDecimal(), LocalDateTime.of(2022, 2, 28, 23, 59)),
                eurUsd to Price(eurUsd, "0.713".toBigDecimal(), "0.71304".toBigDecimal(), LocalDateTime.of(2022, 2, 28, 23, 59))
            ),
            mapOf(
                usdJpy to Price(usdJpy, "113.234".toBigDecimal(), "113.236".toBigDecimal(), LocalDateTime.of(2022, 3, 1, 0, 0)),
            ),
            mapOf(
                usdJpy to Price(usdJpy, "133.234".toBigDecimal(), "133.236".toBigDecimal(), LocalDateTime.of(2022, 3, 1, 0, 1)),
                eurUsd to Price(eurUsd, "0.913".toBigDecimal(), "0.91304".toBigDecimal(), LocalDateTime.of(2022, 3, 1, 0, 1))
            ),
            mapOf(
                eurUsd to Price(eurUsd, "0.510".toBigDecimal(), "0.51004".toBigDecimal(), LocalDateTime.of(2022, 3, 1, 0, 2))
            ),
            mapOf(
                usdJpy to Price(usdJpy, "154.0".toBigDecimal(), "154.002".toBigDecimal(), LocalDateTime.of(2022, 3, 1, 0, 3))
            )
        )
        assertEquals(expected, actual)
    }
}
