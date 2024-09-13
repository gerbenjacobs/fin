package trading212

import (
	"github.com/gerbenjacobs/fin"
)

// Accumulate accumulates all events, does housekeeping, but keeps them as a big list
func Accumulate(splits []fin.Splits, events []TradeEvent) []fin.Aggregate {
	var stocks []fin.Aggregate
	for _, e := range events {
		if !e.IsBuying() && !e.IsSelling() {
			// skip any events that are not buying or selling
			continue
		}

		// deal with buying/selling
		var shareCount = float64(0)
		if e.IsBuying() {
			shareCount = e.ShareCount
		} else {
			shareCount = -e.ShareCount
		}

		// did a stock split happen today
		var avgPrice = e.SharePrice
		for _, split := range splits {
			if split.Symbol == e.TickerSymbol &&
				e.Time.Format("2006-01-02") <= split.Date {
				avgPrice = avgPrice / split.Ratio
				shareCount = shareCount * split.Ratio
			}
		}

		stocks = append(stocks, fin.Aggregate{
			Symbol:     e.TickerSymbol,
			ShareCount: shareCount,
			AvgPrice:   avgPrice,
			LastUpdate: e.Time.Time,
		})
	}
	return stocks
}
