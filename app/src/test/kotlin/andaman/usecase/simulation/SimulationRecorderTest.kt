package andaman.usecase.simulation

import andaman.TestDbAware
import andaman.enum.PositionStatus
import andaman.repository.DbRepository
import andaman.testContext
import andaman.testPosition
import io.mockk.coEvery
import io.mockk.mockk
import io.mockk.verify
import org.junit.After
import org.kodein.di.DI
import org.kodein.di.bindSingleton
import kotlin.test.BeforeTest
import kotlin.test.Test

class SimulationRecorderTest: TestDbAware {

    @BeforeTest
    fun setup() {
        setupTestDb()
    }

    @After
    fun shutdown() {
        shutdownTestDb()
    }

    @Test
    fun testRecord() {
        val repository = mockk<DbRepository>()
        coEvery { repository.insertPositions(any(), any()) } returns Unit
        val di = DI { bindSingleton { repository } }

        val position1 = testPosition()
        val position2 = testPosition().also { it.status = PositionStatus.CLOSED }
        val context = testContext().also { it.currentPositions = listOf(position1, position2) }

        val simulationRecorder = SimulationRecorder(di)
        simulationRecorder.record(context)

        verify { repository.insertPositions(listOf(position2), context.trade.tradeId) }
    }

    @Test
    fun testFinal() {
        val repository = mockk<DbRepository>()
        coEvery { repository.insertPositions(any(), any()) } returns Unit
        val di = DI { bindSingleton { repository } }

        val position1 = testPosition()
        val position2 = testPosition().also { it.status = PositionStatus.CLOSED }
        val context = testContext().also { it.currentPositions = listOf(position1, position2) }

        val simulationRecorder = SimulationRecorder(di)
        simulationRecorder.final(context)

        verify { repository.insertPositions(context.currentPositions, context.trade.tradeId) }
    }
}
