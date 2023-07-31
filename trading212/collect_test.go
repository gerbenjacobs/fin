package trading212

import (
	"reflect"
	"testing"
	"time"
)

// defaultTestDataEvents is a list of TradeEvents that matches testdata/tradign212.csv
// this can be used by several tests
var defaultTestDataEvents = []TradeEvent{
	{Action: "Deposit", Time: DateTime{Time: time.Date(2021, 8, 9, 15, 25, 29, 0, time.UTC)}, Total: 1000, Notes: "Transaction ID: xxx", ID: "d0ca160f-f407-4b9b-bb36-xxx"},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)}, ISIN: "US30303M1027", TickerSymbol: "FB", TickerName: "Meta Platforms", ShareCount: 0.0863914000, SharePrice: 362, ShareCurrency: "USD", ExchangeRate: "1.17437", Total: 26.67, ID: "EOF1", FXFee: 0.04},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)}, ISIN: "US88160R1014", TickerSymbol: "TSLA", TickerName: "Tesla", ShareCount: 0.0766547000, SharePrice: 713.93, ShareCurrency: "USD", ExchangeRate: "1.17437", Total: 46.67, ID: "EOF2", FXFee: 0.07},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2021, 8, 9, 18, 31, 41, 0, time.UTC)}, ISIN: "US5949181045", TickerSymbol: "MSFT", TickerName: "Microsoft", ShareCount: 0.2709950000, SharePrice: 288.53, ShareCurrency: "USD", ExchangeRate: "1.17437", Total: 66.68, ID: "EOF3", FXFee: 0.10},
	{Action: "Market sell", Time: DateTime{Time: time.Date(2021, 8, 30, 13, 30, 3, 0, time.UTC)}, ISIN: "US5949181045", TickerSymbol: "MSFT", TickerName: "Microsoft", ShareCount: 0.2709950000, SharePrice: 301.14, ShareCurrency: "USD", ExchangeRate: "1.17951", Result: 2.61, Total: 69.09, ID: "EOF4", FXFee: 0.10},
	{Action: "Deposit", Time: DateTime{Time: time.Date(2021, 9, 7, 13, 43, 10, 0, time.UTC)}, Total: 1000, Notes: "Transaction ID: xxx", ID: "3e8f5274-1c62-46d6-baf4-xxx"},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2021, 9, 27, 13, 19, 13, 0, time.UTC)}, ISIN: "US02079K1079", TickerSymbol: "ABEC", TickerName: "Google", ShareCount: 0.0041253700, SharePrice: 2424.00, ShareCurrency: "EUR", ExchangeRate: "1.00000", Total: 10.00, ID: "EOF5"},
	{Action: "Dividend (Ordinary)", Time: DateTime{Time: time.Date(2021, 9, 30, 11, 15, 32, 0, time.UTC)}, ISIN: "US5949181045", TickerSymbol: "MSFT", TickerName: "Microsoft", ShareCount: 0.2709950000, SharePrice: 0.48, ShareCurrency: "USD", ExchangeRate: "Not available", Total: 0.11, Tax: 0.02, TaxCurrency: "USD"},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2022, 3, 7, 16, 10, 26, 0, time.UTC)}, ISIN: "FR0000120578", TickerSymbol: "SAN", TickerName: "Sanofi", ShareCount: 0.1117960000, SharePrice: 89.18, ShareCurrency: "EUR", ExchangeRate: "1.00000", Total: 10.00, ID: "EOF6", FRFee: 0.03},
	{Action: "Market buy", Time: DateTime{Time: time.Date(2022, 7, 29, 14, 28, 17, 0, time.UTC)}, ISIN: "US02079K1079", TickerSymbol: "ABEC", TickerName: "Alphabet (Class C)", ShareCount: 2.2887315000, SharePrice: 113.60, ShareCurrency: "EUR", ExchangeRate: "1.00000", Total: 260.00, ID: "EOF7"},
}

func TestCollect(t *testing.T) {
	tests := []struct {
		name           string
		inputDirectory string
		want           []TradeEvent
	}{
		{
			name:           "regular test with our testdata file",
			inputDirectory: "../testdata",
			want:           defaultTestDataEvents,
		},
		{
			name:           "test that GBP is correctly parsed",
			inputDirectory: "../testdata/gbp",
			want:           defaultTestDataEvents,
		},
		{
			name:           "test that multiple files are correctly parsed",
			inputDirectory: "../testdata/multiple",
			want:           defaultTestDataEvents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := Collect(tt.inputDirectory)
			for idx, event := range events {
				if !reflect.DeepEqual(event, tt.want[idx]) {
					t.Errorf("event for index %d is a mismatch \n%#v\n%#v", idx, event, tt.want[idx])
				}
			}
		})
	}
}
