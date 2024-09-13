package trading212

import (
	"reflect"
	"testing"
	"time"

	"github.com/gerbenjacobs/fin"
)

func TestAccumulate(t *testing.T) {
	abecSplitRatio := float64(20)
	splits := []fin.Splits{
		{Symbol: "ABEC", Date: "2022-07-16", Ratio: abecSplitRatio},
	}
	tests := []struct {
		name   string
		events []TradeEvent
		want   []fin.Aggregate
	}{
		{
			name:   "Regular test on our default testdata",
			events: defaultTestDataEvents,
			want: []fin.Aggregate{
				{Symbol: "FB", ShareCount: 0.0863914, AvgPrice: 362, LastUpdate: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)},
				{Symbol: "TSLA", ShareCount: 0.0766547, AvgPrice: 713.93, LastUpdate: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)},
				{Symbol: "MSFT", ShareCount: 0.270995, AvgPrice: 288.53, LastUpdate: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)},
				{Symbol: "MSFT", ShareCount: -0.270995, AvgPrice: 301.14, LastUpdate: time.Date(2021, 8, 30, 13, 30, 3, 0, time.UTC)},
				// this ABEC entry is before the stock split
				{Symbol: "ABEC", ShareCount: 0.0041253700 * abecSplitRatio, AvgPrice: 2424.00 / abecSplitRatio, LastUpdate: time.Date(2021, 9, 27, 13, 19, 13, 0, time.UTC)},
				{Symbol: "SAN", ShareCount: 0.111796, AvgPrice: 89.18, LastUpdate: time.Date(2022, 3, 7, 16, 10, 26, 0, time.UTC)},
				{Symbol: "ABEC", ShareCount: 2.2887315000, AvgPrice: 113.6, LastUpdate: time.Date(2022, 7, 29, 14, 28, 17, 0, time.UTC)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates := Accumulate(splits, tt.events)
			if len(tt.want) != len(aggregates) {
				t.Fatalf("Length mismatch, got %d want %d", len(aggregates), len(tt.want))
			}
			for idx, agg := range aggregates {
				if !reflect.DeepEqual(agg, tt.want[idx]) {
					t.Errorf("Accumulate for %s is a mismatch \n%#v\n%#v", agg.Symbol, agg, tt.want[idx])
				}
			}
		})
	}

}
