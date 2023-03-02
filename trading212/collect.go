package trading212

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gerbenjacobs/fin"
)

// Collect loops through the input directory and parses all CSV files
// it creates trading212 TradeEvents and de-duplicates them
func Collect(inputDirectory string) []TradeEvent {
	// loop through directory and find csv files
	maxDepth := strings.Count(filepath.FromSlash(inputDirectory), string(os.PathSeparator)) + 1
	var fileNames []string
	err := filepath.WalkDir(inputDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.Count(path, string(os.PathSeparator)) > maxDepth {
			// Ignore directories and subdirectories
			return nil
		}
		if filepath.Ext(path) == ".csv" {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to read directory %#v: %v", inputDirectory, err)
	}

	// parse events and skip duplicates
	var events = make(map[string]TradeEvent)
	var sortedEvents []string
	var localEvents []TradeEvent
	for _, fileName := range fileNames {
		err := fin.ReadCSVFile(fileName, &localEvents)
		if err != nil {
			log.Fatalf("Failed to read trading212 data in %s: %v", fileName, err)
		}

		for _, le := range localEvents {
			// create IDs based on time, so we sort them chronologically too
			var id = fmt.Sprintf("%s-%s", le.Time, le.ID)

			// dividends have no ID, generate one
			if le.IsDividend() {
				id = fmt.Sprintf("%s-%s-%s", le.Time, le.TickerSymbol, le.Action)
			}

			// only add new events
			if _, ok := events[id]; !ok {
				sortedEvents = append(sortedEvents, id)
				events[id] = le
			}
		}
	}

	// sort and collate events
	sort.Strings(sortedEvents)
	var finalEvents []TradeEvent
	for _, id := range sortedEvents {
		finalEvents = append(finalEvents, events[id])
	}

	return finalEvents
}
