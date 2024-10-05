package fin

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfig = "config.yaml"

type Config struct {
	Input    string            `yaml:"input"`
	Output   string            `yaml:"output"`
	Format   string            `yaml:"format"`
	SkipPies bool              `yaml:"skip-pies"`
	PieOnly  string            `yaml:"pie-only"`
	Splits   []Splits          `yaml:"splits"`
	Symbols  map[string]string `yaml:"symbols"`
	Renames  map[string]string `yaml:"renames"`
	Pies     []struct {
		Name    string   `yaml:"name"`
		Symbols []string `yaml:"symbols"`
	} `yaml:"pies"`
}

type Splits struct {
	Symbol string  `yaml:"symbol"`
	Date   string  `yaml:"date"`
	Ratio  float64 `yaml:"ratio"`
}

// GetOrCreateConfig checks for a config file and otherwise creates one.
func GetOrCreateConfig() Config {
	cfg, err := GetConfig(defaultConfig)
	switch {
	case os.IsNotExist(err):
		createConfig()
		cfg, _ = GetConfig(defaultConfig)
	case err != nil:
		log.Fatalf("Failed to read config file: %v", err)
	}

	return cfg
}

// GetConfig reads a config file from disk.
func GetConfig(file string) (Config, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// createConfig creates a default config file.
func createConfig() {
	data := `---
# Required config
input: data # folder where your Trading212 CSVs are stored
output: aggregated_quotes # name of output file (prefix)
format: aggregate # aggregate or yahoo

# Optional config
skip-pies: true # skip splitting by pies (default: false)
pie-only: "" # only generate this pie (default: "")

# Splits is a list of split events relevant to your portfolio
# this is needed to calculate the total stock count
splits:
  - symbol: ABEC
    date: 2022-07-16
    ratio: 20 # for reverse splits, use a decimal ratio

# Symbols is a list of conversions to take Trading212 symbols
# and convert them to the symbols used by Yahoo portfolios
symbols:
  RIO: RIO.L
  SAN: SAN.PA

# Pies allows you split your aggregation into multiple CSVs
# uncomment to use
#pies:
#  - name: Growth
#    symbols:
#      - GOOG
#      - AMZN
#  - name: Dividend
#    symbols:
#      - PEP
#      - JNJ
`
	err := os.WriteFile(defaultConfig, []byte(data), 0644)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
}

// PieNames returns the available pie names.
func (c Config) PieNames() []string {
	var names []string
	for _, p := range c.Pies {
		names = append(names, p.Name)
	}
	return names
}
