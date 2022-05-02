package andaman.usecase

/**
 * レコーダー
 */
interface Recorder {
    fun record(context: Context)
    fun final(context: Context)
}
