package andaman.usecase.simulation

import andaman.enum.PositionStatus
import andaman.repository.DbRepository
import andaman.usecase.Context
import andaman.usecase.Recorder
import org.jetbrains.exposed.sql.transactions.transaction
import org.kodein.di.DI
import org.kodein.di.DIAware
import org.kodein.di.instance

class SimulationRecorder(override val di: DI): Recorder, DIAware {
    private val repository by instance<DbRepository>()

    override fun record(context: Context) {
        val positions = context.currentPositions.filter { it.status == PositionStatus.CLOSED }
        transaction {
            repository.insertPositions(positions, context.trade.tradeId)
        }
    }

    override fun final(context: Context) {
        transaction {
            repository.insertPositions(context.currentPositions, context.trade.tradeId)
        }
    }
}
