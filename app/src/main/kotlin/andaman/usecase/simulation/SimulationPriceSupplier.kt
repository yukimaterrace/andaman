package andaman.usecase.simulation

import andaman.enum.CurrencyPair
import andaman.usecase.*
import org.kodein.di.DI
import org.kodein.di.DIAware
import org.kodein.di.instance
import java.math.BigDecimal
import java.nio.file.Path
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.Path

/**
 * シミュレーション用のPriceSupplier
 * （価格データは所定のCSVフォーマットを使用することを前提としています）
 */
class SimulationPriceSupplier(
    private val start: LocalDateTime,
    private val end: LocalDateTime,
    private val currencyPairs: Set<CurrencyPair>,
    private val filePath: String,
    override val di: DI,
): PriceSupplier, DIAware {

    private val productsHolder by instance<ProductsHolder>()

    override fun iterator(): Iterator<CurrencyPairMap> {
        val cps = currencyPairs.map { cp ->
            val datePriceMap = cp.priceFilePaths(start, end, filePath).flatMap { csp ->
                csp.second.toFile().readLines().map { s ->
                    val product = productsHolder.product(csp.first) ?: throw RuntimeException()
                    s.currencyPairPrice(product)
                }
            }.fold(mapOf<LocalDateTime, CurrencyPairPrice>()) { acc, csp ->
                acc.plus(Pair(csp.second.at, csp))
            }
            makeTimeList().map { datePriceMap[it] }
        }
        val priceMaps = timeRange().map { index ->
            cps.fold<List<CurrencyPairPrice?>, CurrencyPairMap>(mapOf()) { acc, csp ->
                val cs = csp[index.toInt()]
                cs?.let { acc.plus(Pair(it.first, it.second)) } ?: acc
            }
        }
        return priceMaps.iterator()
    }

    private fun makeTimeList(): List<LocalDateTime> = timeRange().map { start.plusMinutes(it) }
    private fun timeRange(): LongRange = (0..ChronoUnit.MINUTES.between(start, end))
}

internal typealias CurrencyPairPath = Pair<CurrencyPair, Path>
internal typealias CurrencyPairPrice = Pair<CurrencyPair, Price>

/**
 * 通貨ペアからファイルパスを生成します
 */
internal fun CurrencyPair.priceFilePaths(
    start: LocalDateTime,
    end: LocalDateTime,
    filePath: String
): List<CurrencyPairPath> {
    val startMonth = start.withDayOfMonth(1)
    val numMonths = (end.year - start.year) * 12 + (end.monthValue - start.monthValue) + 1
    return (1..numMonths).map { startMonth.plusMonths(it.toLong() - 1) }.map {
        this to Path(filePath, "${this.fileName()}_${it.year}_${"%02d".format(it.monthValue)}.csv")
    }
}

/**
 * CSVの一行からCurrencyPairPriceを生成します
 */
internal fun String.currencyPairPrice(product: Product): CurrencyPairPrice {
    val source = this.split(",")
    val at = LocalDateTime.parse(
        "%s.%s".format(source[0], source[1]),
        DateTimeFormatter.ofPattern("yyyy.MM.dd.HH:mm")
    )
    val bid = source[5].toBigDecimal()
    val ask = bid + product.pipsUnit * product.simulationSpreadPips
    return product.currencyPair to Price(product.currencyPair, bid, ask, at)
}
