package yahoo

import "time"

const (
	dateFormat = "2006/01/02"
	timeFormat = "15:04 MST"
)

type StockData struct {
	Symbol        string  `csv:"Symbol"`
	SharePrice    float64 `csv:"Current Price"`
	Date          Date    `csv:"Date"`
	Time          Time    `csv:"Time"`
	Change        float64 `csv:"Change"`
	Open          float64 `csv:"Open"`
	High          float64 `csv:"High"`
	Low           float64 `csv:"Low"`
	Volume        int64   `csv:"Volume"`
	TradeDate     string  `csv:"Trade Date"`
	PurchasePrice float64 `csv:"Purchase Price"`
	Quantity      float64 `csv:"Quantity"`
	Commission    string  `csv:"Commission"`
	HighLimit     float64 `csv:"High Limit"`
	LowLimit      float64 `csv:"Low Limit"`
	Comment       string  `csv:"Comment"`
}

type Date struct {
	time.Time
}

func (date *Date) MarshalCSV() (string, error) {
	return date.Time.Format(dateFormat), nil
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(dateFormat, csv)
	return err
}

type Time struct {
	time.Time
}

func (t *Time) MarshalCSV() (string, error) {
	return t.Time.Format(timeFormat), nil
}

func (t *Time) UnmarshalCSV(csv string) (err error) {
	t.Time, err = time.Parse(timeFormat, csv)
	return err
}
