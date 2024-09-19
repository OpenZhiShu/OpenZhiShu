package config

import (
	"OpenZhiShu/pkg/elements"
	"encoding/json"
	"html/template"
	"os"
)

type Config struct {
	HomepageConfig StaticConfig  `json:"homepage"`
	DrawingConfig  StaticConfig  `json:"drawing"`
	ResultConfig   DynamicConfig `json:"result"`
}

type StaticConfig struct {
	BodyColor string             `json:"body_color"`
	Ratio     template.CSS       `json:"ratio"`
	Elements  []elements.Element `json:"elements"`
}

type DynamicConfig struct {
	BodyColor string             `json:"body_color"`
	Ratio     template.CSS       `json:"ratio"`
	Elements  []elements.Element `json:"elements"`
}

func (c Config) Verify() error {
	if err := c.HomepageConfig.Verify(); err != nil {
		return err
	}
	if err := c.DrawingConfig.Verify(); err != nil {
		return err
	}
	if err := c.ResultConfig.Verify(); err != nil {
		return err
	}
	return nil
}

func (s StaticConfig) Verify() error {
	for i := range s.Elements {
		if err := s.Elements[i].Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (d DynamicConfig) Verify() error {
	for i := range d.Elements {
		if err := d.Elements[i].Verify(); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig(filepath string) (Config, error) {
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		return Config{}, err
	}

	if err = cfg.Verify(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
