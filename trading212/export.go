package trading212

import (
	"strings"
	"time"
)

const format = "2006-01-02 15:04:05"

func (e *TradeEvent) IsBuying() bool {
	return strings.Contains(e.Action, "buy")
}

func (e *TradeEvent) IsSelling() bool {
	return strings.Contains(e.Action, "sell")
}

func (e *TradeEvent) IsDividend() bool {
	return strings.Contains(e.Action, "Dividend")
}

func (e *TradeEvent) Fees() float64 {
	return e.FXFee + e.FRFee + e.StampDuty + e.StampDutyTax + e.FinraFee
}

type TradeEvent struct {
	Action        string
	Time          DateTime `csv:"Time"`
	ISIN          string   `csv:"ISIN"`
	TickerSymbol  string   `csv:"Ticker"`
	TickerName    string   `csv:"Name"`
	ShareCount    float64  `csv:"No. of shares,omitempty"`
	SharePrice    float64  `csv:"Price / share,omitempty"`
	ShareCurrency string   `csv:"Currency (Price / share)"`
	ExchangeRate  string   `csv:"Exchange rate,omitempty"`
	Result        float64  `csv:"Result,omitempty"` // gain or loss
	Total         float64  `csv:"Total,omitempty"`  // total money gained
	Tax           float64  `csv:"Withholding tax,omitempty"`
	TaxCurrency   string   `csv:"Currency (Withholding tax)"`
	Deposit       float64  `csv:"Charge amount,omitempty"` // amount of money when depositing
	DepositFee    float64  `csv:"Deposit fee,omitempty"`   // fee paid to trading212
	StampDuty     float64  `csv:"Stamp duty,omitempty"`
	StampDutyTax  float64  `csv:"Stamp duty reserve tax,omitempty"`
	Notes         string   `csv:"Notes"`
	ID            string   `csv:"ID"`
	FXFee         float64  `csv:"Currency conversion fee,omitempty"` // foreign exchange fee
	FRFee         float64  `csv:"French transaction tax,omitempty"`
	FinraFee      float64  `csv:"Finra fee,omitempty"`
}

type DateTime struct {
	time.Time
}

func (dt *DateTime) UnmarshalCSV(csv string) (err error) {
	dt.Time, err = time.Parse(format, csv)
	return
}

func (dt *DateTime) MarshalCSV() (string, error) {
	return dt.Format(format), nil
}
