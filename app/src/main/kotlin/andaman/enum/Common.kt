package andaman.enum

/**
 * 通貨ペア
 */
enum class CurrencyPair {
    UsdJpy,
    EurJpy,
    GbpJpy,
    EurUsd,
    GbpUsd,
    EurGbp
}

/**
 * ポジションステータス
 */
enum class PositionStatus {
    OPEN,
    CLOSED
}

/**
 * 売買タイプ
 */
enum class BuySellType {
    BUY,
    SELL
}
