package techan

import (
	"github.com/sdcoffey/big"
)

type adxIndicator struct {
	series *TimeSeries
	window int
	// unstablePeriod int
}

func NewADXAndDIIndicator(series *TimeSeries, window int, unstablePeriod int) TripleIndicator {
	return &adxIndicator{
		series: series,
		window: window,
	}
}

func (adxi *adxIndicator) Calculate(index int) (big.Decimal, big.Decimal, big.Decimal) {
	if index < adxi.window {
		return big.ZERO, big.ZERO, big.ZERO
	}

	minusDm := NewMinusDMIndicator(adxi.series, adxi.window)
	plusDm := NewPlusDMIndicator(adxi.series, adxi.window)
	plusDi := NewPlusDiIndicator(plusDm, adxi.series, adxi.window)
	minusDi := NewMinusDiIndicator(minusDm, adxi.series, adxi.window)

	dx := NewDxIndicator(plusDi, minusDi, adxi.window)

	adx := NewAdxRawIndicator(dx, adxi.window)
	adx_smooth := NewEMAIndicator(adx, adxi.window)

	return adx_smooth.Calculate(index), plusDi.Calculate(index), minusDi.Calculate(index)
}

type minusDiIndicator struct {
	minusDm     Indicator
	series      *TimeSeries
	window      int
	resultCache resultCache
}

func NewMinusDiIndicator(minusDm Indicator, series *TimeSeries, window int) Indicator {
	return &minusDiIndicator{
		minusDm: minusDm,
		series:  series,
		window:  window,
	}
}

func (mdi *minusDiIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(mdi, index, func(i int) big.Decimal {
		atr := NewSimpleMovingAverage(NewTrueRangeIndicator(mdi.series), mdi.window).Calculate(index)
		if atr.EQ(big.ZERO) {
			return big.ZERO
		}
		return NewEMAIndicator(mdi.minusDm, mdi.window).Calculate(index).Div(atr).Mul(big.NewFromInt(100))
	}); cachedValue != nil {
		return *cachedValue
	}
	atr := NewSimpleMovingAverage(NewTrueRangeIndicator(mdi.series), mdi.window).Calculate(index)
	if atr.EQ(big.ZERO) {
		return big.ZERO
	}
	minusDi := NewEMAIndicator(mdi.minusDm, mdi.window).Calculate(index).Div(atr).Mul(big.NewFromInt(100))

	cacheResult(mdi, index, minusDi)

	return minusDi
}

func (mdi minusDiIndicator) cache() resultCache { return mdi.resultCache }

func (mdi *minusDiIndicator) setCache(newCache resultCache) {
	mdi.resultCache = newCache
}

func (mdi minusDiIndicator) windowSize() int { return mdi.window }

type plusDiIndicator struct {
	plusDm      Indicator
	series      *TimeSeries
	window      int
	resultCache resultCache
}

func NewPlusDiIndicator(plusDm Indicator, series *TimeSeries, window int) Indicator {
	return &plusDiIndicator{
		plusDm: plusDm,
		series: series,
		window: window,
	}
}

func (pdi *plusDiIndicator) Calculate(index int) big.Decimal {
	if cachedValue := returnIfCached(pdi, index, func(i int) big.Decimal {
		atr := NewSimpleMovingAverage(NewTrueRangeIndicator(pdi.series), pdi.window).Calculate(index)
		if atr.EQ(big.ZERO) {
			return big.ZERO
		}
		return NewEMAIndicator(pdi.plusDm, pdi.window).Calculate(index).Div(atr).Mul(big.NewFromInt(100))
	}); cachedValue != nil {
		return *cachedValue
	}
	atr := NewSimpleMovingAverage(NewTrueRangeIndicator(pdi.series), pdi.window).Calculate(index)
	if atr.EQ(big.ZERO) {
		return big.ZERO
	}
	plusDi := NewEMAIndicator(pdi.plusDm, pdi.window).Calculate(index).Div(atr).Mul(big.NewFromInt(100))

	cacheResult(pdi, index, plusDi)

	return plusDi
}

func (pdi plusDiIndicator) cache() resultCache { return pdi.resultCache }

func (pdi *plusDiIndicator) setCache(newCache resultCache) {
	pdi.resultCache = newCache
}

func (pdi plusDiIndicator) windowSize() int { return pdi.window }

type dxIndicator struct {
	plusDi  Indicator
	minusDi Indicator
	window  int
}

func NewDxIndicator(plusDi Indicator, minusDi Indicator, window int) Indicator {
	return &dxIndicator{
		plusDi:  plusDi,
		minusDi: minusDi,
		window:  window,
	}
}

func (dxi *dxIndicator) Calculate(index int) big.Decimal {
	if index < dxi.window {
		return big.ZERO
	}
	plusDi := dxi.plusDi.Calculate(index)
	minusDi := dxi.minusDi.Calculate(index)
	if plusDi.EQ(big.ZERO) && minusDi.EQ(big.ZERO) {
		return big.ZERO
	}
	return (plusDi.Sub(minusDi).Abs()).Div((plusDi.Add(minusDi).Abs())).Mul(big.NewFromInt(100))
}

type adxRawIndicator struct {
	dxi    Indicator
	window int
}

func NewAdxRawIndicator(dxi Indicator, window int) Indicator {
	return &adxRawIndicator{
		dxi:    dxi,
		window: window,
	}
}

func (adxri *adxRawIndicator) Calculate(index int) big.Decimal {
	return (adxri.dxi.Calculate(index - 1).Mul(big.NewFromInt(adxri.window - 1)).Add(adxri.dxi.Calculate(index))).Div(big.NewFromInt(adxri.window))
}
