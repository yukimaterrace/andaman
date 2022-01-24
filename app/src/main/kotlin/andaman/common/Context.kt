package andaman.common

import andaman.price.CurrencySymbol
import andaman.price.Price

class Context {
    var currentPrices: Map<CurrencySymbol, Price> = emptyMap()
}
