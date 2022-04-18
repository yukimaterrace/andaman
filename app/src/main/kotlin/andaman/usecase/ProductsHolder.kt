package andaman.usecase

import andaman.enum.CurrencyPair
import com.charleskorn.kaml.Yaml
import kotlinx.serialization.*
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.descriptors.serialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import java.math.BigDecimal
import kotlin.io.path.Path

class ProductsHolder(private val sourcePath: String) {
    private var productMap: Map<CurrencyPair, Product> = emptyMap()

    /**
     * 所定のYamlファイルから商品データをロードします
     */
    fun load() {
        val data = Path(sourcePath).toFile().readText()
        val products = Yaml.default.decodeFromString<Products>(data)
        productMap = products.products.associateBy { it.currencyPair }
    }

    /**
     * 通貨ペアからプロダクトを取得します
     */
    fun product(currencyPair: CurrencyPair) = productMap[currencyPair]
}

/**
 * 商品情報
 */
@Serializable
data class Product(
    val currencyPair: CurrencyPair,
    @Serializable(with = BigDecimalSerializer::class)
    val pipsUnit: BigDecimal,
    @Serializable(with = BigDecimalSerializer::class)
    val simulationSpreadPips: BigDecimal
)

@Serializable
private data class Products(val products: List<Product>)

/**
 * BigDecimalのシリアライザ
 */
object BigDecimalSerializer: KSerializer<BigDecimal> {
    override fun deserialize(decoder: Decoder): BigDecimal {
        return decoder.decodeString().toBigDecimal()
    }

    override fun serialize(encoder: Encoder, value: BigDecimal) {
        encoder.encodeString(value.toPlainString())
    }

    override val descriptor: SerialDescriptor
        get() = PrimitiveSerialDescriptor("BigDecimal", PrimitiveKind.STRING)
}
