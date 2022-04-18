/*
 * This Kotlin source file was generated by the Gradle 'init' task.
 */
package andaman

import andaman.usecase.Context
import andaman.usecase.ProductsHolder
import andaman.usecase.Trade
import andaman.usecase.User
import org.kodein.di.DI
import org.kodein.di.bindSingleton
import java.util.*

fun testUser() = User(accountId = 0, name = "")
fun testTrade() = Trade(tradeId = UUID.randomUUID())
fun testContext() = Context(testUser(), testTrade())

fun testDI() = DI {
    bindSingleton { ProductsHolder("resource/test/test_products.yaml").also { it.load() } }
}
