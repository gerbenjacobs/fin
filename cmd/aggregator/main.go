package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/fin"
	"github.com/gerbenjacobs/fin/trading212"
	"github.com/gerbenjacobs/fin/yahoo"
)

const outputNone = "none"

var log *logrus.Logger

func main() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	// deal with config
	fConfig := flag.String("config", "", "YAML file with configuration for this run")
	skipPies := flag.Bool("skip-pies", false, "Skip pie separation during output")
	pieOnly := flag.String("pie-only", "", "Only output specific pie output")
	flag.Parse()

	var cfg fin.Config
	if fConfig == nil || *fConfig == "" {
		cfg = fin.GetOrCreateConfig()
	} else {
		var err error
		cfg, err = fin.GetConfig(*fConfig)
		if err != nil {
			log.Fatalf("Failed to read config file %s: %v", *fConfig, err)
		}
	}

	// overwrite config with flags
	if skipPies != nil && *skipPies {
		cfg.SkipPies = *skipPies
	}
	if pieOnly != nil && *pieOnly != "" {
		cfg.PieOnly = *pieOnly
	}

	// check if this pie exists
	if cfg.PieOnly != "" {
		var found bool
		for _, p := range cfg.PieNames() {
			if p == cfg.PieOnly {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Pie %#v does not exist, we found %s", cfg.PieOnly, cfg.PieNames())
		}
	}

	log.WithFields(logrus.Fields{
		"input":     cfg.Input,
		"output":    cfg.Output,
		"format":    cfg.Format,
		"skip-pies": cfg.SkipPies,
		"pie-only":  cfg.PieOnly,
		"pies":      len(cfg.Pies),
		"splits":    len(cfg.Splits),
		"symbols":   len(cfg.Symbols),
		"renames":   len(cfg.Renames),
	}).Info("Starting process.")

	// loop through directory and find csv files
	events := trading212.Collect(cfg.Input)

	// aggregate events via Trading212 algorithm
	stocks, totals := trading212.Aggregate(cfg.Splits, cfg.Renames, events)

	log.WithFields(logrus.Fields{
		"deposits":            totals.Deposits,
		"invested":            totals.Invested,
		"realized":            totals.Realized,
		"realized-with-costs": ceilFloat(totals.Realized-totals.Fees-totals.Taxes, 2),
		"dividends":           totals.Dividends,
		"fees":                totals.Fees,
		"taxes":               totals.Taxes,
		"cash":                totals.Cash,
		"interest":            totals.Interest,
		"withdrawals":         totals.Withdrawals,
	}).Info("Completed aggregation.")

	// write output
	switch cfg.Format {
	case "aggregate":
		var output = make(map[string][]fin.Aggregate)
		for _, s := range stocks {
			outputName := determineOutputName(cfg, s.Symbol)

			if cfg.PieOnly != "" && cfg.PieOnly != outputName {
				// pieOnly is set, skip all entries that don't match
				continue
			}
			if cfg.SkipPies {
				outputName = outputNone
			}

			// add output
			output[outputName] = append(output[outputName], s)
		}
		for on, o := range output {
			fn := writeOutputJSON(cfg, on, o)
			log.Printf("Written %d entries to %s.", len(o), fn)
		}
	case "yahoo":
		var output = make(map[string][]yahoo.StockData)
		for _, a := range stocks {
			if a.ShareCount == 0 {
				// skip entries with no quantity
				continue
			}

			outputName := determineOutputName(cfg, a.Symbol) // use original symbol

			if cfg.PieOnly != "" && cfg.PieOnly != outputName {
				// pieOnly is set, skip all entries that don't match
				continue
			}
			if cfg.SkipPies {
				outputName = outputNone
			}

			sym := a.Symbol
			if v, ok := cfg.Symbols[sym]; ok {
				sym = v
			}

			// add output
			output[outputName] = append(output[outputName], yahoo.StockData{
				Symbol:        sym,
				Quantity:      a.ShareCount,
				PurchasePrice: a.AvgPrice,
				TradeDate:     a.LastUpdate.Format("20060102"),
			})
		}

		for on, o := range output {
			fn := writeOutput(cfg, on, o)
			log.Printf("Written %d entries to %s using %d symbol replacements.", len(o), fn, len(cfg.Symbols))
		}
	default:
		log.Fatalf("Unknown output format: %s", cfg.Format)
	}
}

// determineOutputName determines the pie name for a given stock
// or 'none' if none matches
func determineOutputName(cfg fin.Config, s string) string {
	outputName := outputNone
	stockFound := false
	for _, p := range cfg.Pies {
		for _, t := range p.Symbols {
			if t == s {
				outputName = p.Name
				stockFound = true
				break
			}
		}
		if stockFound {
			break
		}
	}

	return outputName
}

// writeOutput creates the output file name and writes the data to it
// then returns the generated file name
func writeOutput(cfg fin.Config, outputName string, output interface{}) string {
	fn := fmt.Sprintf("%s_%s.csv", cfg.Output, strings.ToLower(outputName))
	if outputName == outputNone {
		fn = cfg.Output + ".csv"
	}
	if err := fin.WriteCSVFile(fn, output); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	return fn
}

func writeOutputJSON(cfg fin.Config, outputName string, output interface{}) string {
	fn := fmt.Sprintf("%s_%s.json", cfg.Output, strings.ToLower(outputName))
	if outputName == outputNone {
		fn = cfg.Output + ".json"
	}
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	err = enc.Encode(output)
	if err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	return fn
}

func ceilFloat(f float64, precision int) float64 {
	d := math.Pow(10, float64(precision))
	return math.Ceil(f*d) / d
}
