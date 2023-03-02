package fin

import "time"

type Aggregate struct {
	Symbol         string
	Name           string
	ShareCount     float64
	AvgPrice       float64
	PriceCurrency  string    // currency that goes with avg price; EUR, USD, GBX
	ShareCost      float64   // amount of money spent on this stock, based on the stock's currency. Only used for avg price
	ShareCostLocal float64   // amount of money spent on this stock, based on the portfolio's currency
	ShareResult    float64   // amount of money received from selling
	TotalDividend  float64   // amount of money received from dividends
	Fees           float64   // FXFee, FRFee, StampDutyTax
	Final          float64   // Result + Dividends, minus Cost and Fees.
	LastUpdate     time.Time // used internally to deal with splits
}

type Totals struct {
	Deposits    float64 // the money you deposited
	DepositFees float64 // the fees for depositing
	Invested    float64 // the money you have invested, minus fees
	Realized    float64 // gains you have realized by selling
	Dividends   float64 // amount of money you received from dividends
	Fees        float64 // fees you paid
	Cash        float64 // cash left in your portfolio
	Taxes       float64 // taxes withheld from dividends
}
