package trading212

import (
	"math"
	"sort"

	"github.com/gerbenjacobs/fin"
)

// Aggregate takes a map of events and aggregates them into a map of stocks and totals,
// based on the Trading212 algorithm, along with stock splits.
func Aggregate(splits []fin.Splits, renames map[string]string, events []TradeEvent) ([]fin.Aggregate, fin.Totals) {
	var stocks = make(map[string]fin.Aggregate)
	var stockNames []string
	var totals fin.Totals
	for _, e := range events {
		// skip every event we don't deal with
		if e.IsSkippable() {
			continue
		}
		// handle deposits or additions
		if e.Action == "Deposit" || e.Action == "Spending cashback" {
			totals.Deposits += e.Total
			totals.Withdrawals -= e.DepositFee
			continue
		}
		// handle interest
		if e.IsInterest() {
			totals.Interest += e.Total
			continue
		}
		// handle money withdrawl
		if e.IsMoneyWithdrawal() {
			// we subtract the total, because it's stored as a negative number
			totals.Withdrawals -= e.Total
			continue
		}

		// if no action matches, but our symbol is empty, we continue too
		if e.TickerSymbol == "" {
			continue
		}

		// handle renamed stock symbols
		symbol := e.TickerSymbol
		if rn, ok := renames[symbol]; ok {
			symbol = rn
		}

		// create entry if it doesn't exist
		if _, ok := stocks[symbol]; !ok {
			stocks[symbol] = fin.Aggregate{
				Symbol: symbol,
			}
			stockNames = append(stockNames, symbol)
		}

		// calculate changes
		a := stocks[symbol]

		// did a stock split happen today
		for _, split := range splits {
			if split.Symbol == symbol &&
				split.Date > a.LastUpdate.Format("2006-01-02") && split.Date <= e.Time.Format("2006-01-02") {
				a.ShareCount = a.ShareCount * split.Ratio
			}
		}
		if e.IsBuying() {
			a.ShareCount += e.ShareCount
			a.ShareCost += ceilFloat(e.ShareCount*e.SharePrice, 3)
			a.ShareCostLocal += e.Total
			totals.Invested += e.Total - e.Fees()
		}
		if e.IsSelling() {
			a.ShareCount -= e.ShareCount
			a.ShareCost -= e.Total
			a.ShareCostLocal -= e.Total
			a.ShareResult += e.Result
			totals.Realized += e.Result
			totals.Invested -= e.Total - e.Result + e.Fees()
		}
		if e.IsDividend() {
			a.TotalDividend += e.Total
			totals.Dividends += e.Total
		}

		// calculate all fees and update
		a.ShareCount = floorFloat(a.ShareCount, 6)
		a.Fees += e.Fees()
		totals.Fees += e.Fees()
		totals.Taxes += e.Tax

		// update totals
		if floorFloat(a.ShareCount, 4) > 0 {
			// if it's practically zero, reset it (float comparison issues)
			a.AvgPrice = a.ShareCost / a.ShareCount
		} else {
			// during this event everything was sold
			// update certain fields, so they can be re-used again
			a.ShareCount = 0
			a.ShareCost = 0
			a.ShareCostLocal = 0
			a.AvgPrice = 0
		}

		a.Final = a.ShareResult + a.TotalDividend - a.Fees
		if e.Time.Time.After(a.LastUpdate) {
			a.LastUpdate = e.Time.Time
		}
		a.Name = e.TickerName
		a.PriceCurrency = e.ShareCurrency

		// update entry in map with the latest information
		stocks[a.Symbol] = a
	}

	// calculate splits for untouched stocks
	for _, split := range splits {
		s := stocks[split.Symbol]
		if split.Date > s.LastUpdate.Format("2006-01-02") {
			s.ShareCount = s.ShareCount * split.Ratio
		}
		stocks[split.Symbol] = s
	}

	// calculate cash left over in portfolio
	moneyGained := totals.Deposits + totals.Realized + totals.Dividends
	moneySpent := totals.Invested + totals.Fees + totals.Withdrawals
	totals.Cash = moneyGained - moneySpent

	// format money values to 2 decimals
	for s, stock := range stocks {
		stock.AvgPrice = floorFloat(stock.AvgPrice, 2)
		stock.ShareCost = floorFloat(stock.ShareCost, 2)
		stock.ShareCostLocal = floorFloat(stock.ShareCostLocal, 2)
		stock.ShareResult = floorFloat(stock.ShareResult, 2)
		stock.TotalDividend = floorFloat(stock.TotalDividend, 2)
		stock.Fees = floorFloat(stock.Fees, 2)
		stock.Final = floorFloat(stock.Final, 2)
		stocks[s] = stock
	}
	totals.Deposits = floorFloat(totals.Deposits, 2)
	totals.Invested = floorFloat(totals.Invested, 2)
	totals.Realized = floorFloat(totals.Realized, 2)
	totals.Dividends = floorFloat(totals.Dividends, 2)
	totals.Fees = floorFloat(totals.Fees, 2)
	totals.Cash = floorFloat(totals.Cash, 2)
	totals.Taxes = floorFloat(totals.Taxes, 2)
	totals.Withdrawals = floorFloat(totals.Withdrawals, 2)

	// sort and collate aggregates
	sort.Strings(stockNames)
	var aggregates []fin.Aggregate
	for _, id := range stockNames {
		aggregates = append(aggregates, stocks[id])
	}

	return aggregates, totals
}

func floorFloat(f float64, precision int) float64 {
	d := math.Pow(10, float64(precision))
	return math.Floor(f*d) / d
}

func ceilFloat(f float64, precision int) float64 {
	d := math.Pow(10, float64(precision))
	return math.Ceil(f*d) / d
}
