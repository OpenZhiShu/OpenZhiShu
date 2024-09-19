package config

import (
	"OpenZhiShu/pkg/elements"
	"encoding/json"
	"html/template"
	"os"
)

type Config struct {
	BodyColor string             `json:"body_color"`
	Ratio     template.CSS       `json:"ratio"`
	Elements  []elements.Element `json:"elements"`
}

func (c Config) Verify() error {
	for i := range c.Elements {
		if err := c.Elements[i].Verify(); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig[T interface{ Verify() error }](filepath string) (T, error) {
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		var t T
		return t, err
	}

	var cfg T
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		var t T
		return t, err
	}

	if err = cfg.Verify(); err != nil {
		var t T
		return t, err
	}

	return cfg, nil
}
