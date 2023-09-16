package techan

import "github.com/sdcoffey/big"

type minusDmIndicator struct {
	series *TimeSeries
	window int
}

func NewMinusDMIndicator(series *TimeSeries, window int) Indicator {
	return &minusDmIndicator{
		series: series,
		window: window,
	}
}

func (dmi *minusDmIndicator) Calculate(index int) big.Decimal {
	i := index - 2*dmi.window
	if i < 1 {
		i = 1
	}
	minusDm := big.ZERO
	plusDm := big.ZERO
	for ; i <= index; i++ {
		diffP := dmi.series.Candles[i].MaxPrice.Sub(dmi.series.Candles[i-1].MaxPrice)
		diffM := dmi.series.Candles[index-1].MinPrice.Sub(dmi.series.Candles[index].MinPrice)
		if diffM.GT(big.ZERO) && diffP.LT(diffM) {
			minusDm = minusDm.Add(diffM)
		} else if diffP.GT(big.ZERO) && diffP.GT(diffM) {
			plusDm = plusDm.Add(diffP)
		}
	}
	return minusDm
}

type plusDmIndicator struct {
	series *TimeSeries
	window int
}

func NewPlusDMIndicator(series *TimeSeries, window int) Indicator {
	return &plusDmIndicator{
		series: series,
		window: window,
	}
}

func (dmi *plusDmIndicator) Calculate(index int) big.Decimal {
	i := index - 2*dmi.window
	if i < 1 {
		i = 1
	}
	minusDm := big.ZERO
	plusDm := big.ZERO
	for ; i <= index; i++ {
		diffP := dmi.series.Candles[i].MaxPrice.Sub(dmi.series.Candles[i-1].MaxPrice)
		diffM := dmi.series.Candles[index-1].MinPrice.Sub(dmi.series.Candles[index].MinPrice)
		if diffM.GT(big.ZERO) && diffP.LT(diffM) {
			minusDm = minusDm.Add(diffM)
		} else if diffP.GT(big.ZERO) && diffP.GT(diffM) {
			plusDm = plusDm.Add(diffP)
		}
	}
	return plusDm
}
