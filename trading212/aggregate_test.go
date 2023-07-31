package trading212

import (
	"reflect"
	"testing"
	"time"

	"github.com/gerbenjacobs/fin"
)

func TestAggregate(t *testing.T) {
	splits := []fin.Splits{
		{Symbol: "ABEC", Date: "2022-07-16", Ratio: 20},
	}
	tests := []struct {
		name   string
		events []TradeEvent
		want   []fin.Aggregate
		totals *fin.Totals
	}{
		{
			name:   "Regular test like our testdata",
			events: defaultTestDataEvents,
			want: []fin.Aggregate{
				{Symbol: "ABEC", Name: "Alphabet (Class C)", ShareCount: 2.371231, AvgPrice: 113.86, PriceCurrency: "EUR", ShareCost: 270, ShareCostLocal: 270, ShareResult: 0, TotalDividend: 0, Fees: 0, Final: 0, LastUpdate: time.Date(2022, 7, 29, 14, 28, 17, 0, time.UTC)},
				{Symbol: "FB", Name: "Meta Platforms", ShareCount: 0.086391, AvgPrice: 362, PriceCurrency: "USD", ShareCost: 31.27, ShareCostLocal: 26.67, ShareResult: 0, TotalDividend: 0, Fees: 0.04, Final: -0.04, LastUpdate: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)},
				{Symbol: "MSFT", Name: "Microsoft", ShareCount: 0, AvgPrice: 0, PriceCurrency: "USD", ShareCost: 0, ShareCostLocal: 0, ShareResult: 2.61, TotalDividend: 0.11, Fees: 0.2, Final: 2.51, LastUpdate: time.Date(2021, 9, 30, 11, 15, 32, 0, time.UTC)},
				{Symbol: "SAN", Name: "Sanofi", ShareCount: 0.111796, AvgPrice: 89.18, PriceCurrency: "EUR", ShareCost: 9.97, ShareCostLocal: 10, ShareResult: 0, TotalDividend: 0, Fees: 0.03, Final: -0.03, LastUpdate: time.Date(2022, 3, 7, 16, 10, 26, 0, time.UTC)},
				{Symbol: "TSLA", Name: "Tesla", ShareCount: 0.076654, AvgPrice: 713.94, PriceCurrency: "USD", ShareCost: 54.72, ShareCostLocal: 46.67, ShareResult: 0, TotalDividend: 0, Fees: 0.07, Final: -0.08, LastUpdate: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)},
			},
			totals: &fin.Totals{
				Deposits:  2000,
				Invested:  353.2,
				Realized:  2.61,
				Dividends: 0.11,
				Fees:      0.34,
				Cash:      1649.17,
				Taxes:     0.02,
			},
		},
		{
			name: "Test with a split",
			events: []TradeEvent{
				{Action: "Market buy", Time: DateTime{Time: time.Date(2021, 9, 27, 13, 19, 13, 0, time.UTC)}, TickerSymbol: "ABEC", ShareCount: 0.005, SharePrice: 2000.00, Total: 10.00, ID: "EOF1"},
				{Action: "Market buy", Time: DateTime{Time: time.Date(2022, 9, 27, 13, 19, 13, 0, time.UTC)}, TickerSymbol: "ABEC", ShareCount: 0.125, SharePrice: 80.00, Total: 10.00, ID: "EOF2"},
			},
			want: []fin.Aggregate{
				{Symbol: "ABEC", ShareCount: 0.225, AvgPrice: 88.88, ShareCost: 20, ShareCostLocal: 20, ShareResult: 0, TotalDividend: 0, Fees: 0, Final: 0, LastUpdate: time.Date(2022, 9, 27, 13, 19, 13, 0, time.UTC)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, totals := Aggregate(splits, tt.events)
			for idx, agg := range aggregates {
				if !reflect.DeepEqual(agg, tt.want[idx]) {
					t.Errorf("aggregate for %s is a mismatch \n%#v\n%#v", agg.Symbol, agg, tt.want[idx])
				}
			}

			if tt.totals != nil {
				if !reflect.DeepEqual(totals, *tt.totals) {
					t.Errorf("totals are a mismatch \n%#v\n%#v", totals, tt.totals)
				}
			}
		})
	}
}
