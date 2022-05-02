package andaman.usecase

import andaman.testCloseOrderProposal
import andaman.testContext
import andaman.testOpenOrderProposal
import andaman.usecase.strategy.OrderProposal
import andaman.usecase.strategy.TradeStrategy
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import kotlin.test.Test

class TraderTest {

    @Test
    fun testTrade() {
        val openOrderProposal = testOpenOrderProposal()
        val closeProposal = testCloseOrderProposal()
        val orderProposal = OrderProposal(
            opens = listOf(openOrderProposal),
            closes = listOf(closeProposal)
        )

        val strategy = mockk<TradeStrategy>()
        every { strategy.orderProposal(any()) } returns orderProposal

        val context = testContext(listOf(strategy))
        val broker = mockk<Broker>()
        every { broker.makeMarketOrder(any(), any(), any(), any()) } returns null
        every { broker.closeOrder(any(), any()) } returns null

        val trader = TraderImpl(broker)
        trader.trade(context)

        verify {
            broker.makeMarketOrder(context, openOrderProposal.currencyPair, openOrderProposal.buySellType, openOrderProposal.amount)
        }
        verify {
            broker.closeOrder(context, closeProposal.positionId)
        }
    }
}
