package andaman.price

import java.math.BigDecimal
import java.nio.file.Path
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.Path

class SimulationPriceSupplier(
    private val start: LocalDateTime,
    private val end: LocalDateTime,
    private val symbols: Set<CurrencySymbol>,
    private val filePath: String
): PriceSupplier {

    override fun iterator(): Iterator<PriceMap> {
        val cps = symbols.map { cs ->
            val datePriceMap = cs.priceFilePaths(start, end, filePath).flatMap { csp ->
                csp.second.toFile().readLines().map { s -> s.currencySymbolPrice(csp.first) }
            }.fold(mapOf<LocalDateTime, CurrencySymbolPrice>()) { acc, csp ->
                acc.plus(Pair(csp.second.at, csp))
            }
            makeTimeList().map { datePriceMap[it] }
        }
        val priceMaps = timeRange().map { index ->
            cps.fold<List<CurrencySymbolPrice?>, PriceMap>(mapOf()) { acc, csp ->
                val cs = csp[index.toInt()]
                cs?.let { acc.plus(Pair(it.first, it.second)) } ?: acc
            }
        }
        return priceMaps.iterator()
    }

    private fun makeTimeList(): List<LocalDateTime> = timeRange().map { start.plusMinutes(it) }
    private fun timeRange(): LongRange = (0..ChronoUnit.MINUTES.between(start, end))
}

internal typealias CurrencySymbolPath = Pair<CurrencySymbol, Path>
internal typealias CurrencySymbolPrice = Pair<CurrencySymbol, Price>

internal fun CurrencySymbol.priceFilePaths(
    start: LocalDateTime,
    end: LocalDateTime,
    filePath: String
): List<CurrencySymbolPath> {
    val startMonth = start.withDayOfMonth(1)
    val numMonths = (end.year - start.year) * 12 + (end.monthValue - start.monthValue) + 1
    return (1..numMonths).map { startMonth.plusMonths(it.toLong() - 1) }.map {
        Pair(this, Path(filePath, "${this.fileName()}_${it.year}_${"%02d".format(it.monthValue)}.csv"))
    }
}

internal fun String.currencySymbolPrice(symbol: CurrencySymbol): CurrencySymbolPrice {
    val source = this.split(",")
    val at = LocalDateTime.parse(
        "%s.%s".format(source[0], source[1]),
        DateTimeFormatter.ofPattern("yyyy.MM.dd.HH:mm")
    )
    val value = BigDecimal(source[5])
    return Pair(symbol, Price(symbol, value, at))
}
