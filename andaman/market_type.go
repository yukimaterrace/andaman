package andaman

type (
	instrument interface {
		String() string
	}

	gbpUsd instrument
	eurUsd instrument
	audUsd instrument
	audJpy instrument
	gbpAud instrument
	eurAud instrument
	usdJpy instrument
	gbpJpy instrument
	eurJpy instrument
)

type (
	marketTime interface {
		String()
	}

	s5  marketTime
	s15 marketTime
	m1  marketTime
	m5  marketTime
	m15 marketTime
	h1  marketTime
	h4  marketTime
)

type (
	priceType interface {
		String()
	}

	bid priceType
	ask priceType
)

type (
	orderType interface {
		String()
	}

	buy  orderType
	sell orderType
)
