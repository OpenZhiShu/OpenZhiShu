package config

import (
	"OpenZhiShu/pkg/renderable"
	"encoding/json"
	"html/template"
	"os"
)

type Config struct {
	HomepageConfig HomepageConfig `json:"homepage"`
	DrawingConfig  DrawingConfig  `json:"drawing"`
}

type HomepageConfig struct {
	BodyColor  string                `json:"body_color"`
	Ratio      template.CSS          `json:"ratio"`
	Background renderable.Background `json:"background"`
	Elements   []renderable.Element  `json:"elements"`
}

type DrawingConfig struct {
	BodyColor  string                `json:"body_color"`
	Ratio      string                `json:"ratio"`
	Background renderable.Background `json:"background"`
	Elements   []renderable.Element  `json:"elements"`
}

func (c Config) Verify() error {
	if err := c.HomepageConfig.Verify(); err != nil {
		return err
	}
	if err := c.DrawingConfig.Verify(); err != nil {
		return err
	}
	return nil
}

func (h HomepageConfig) Verify() error {
	if err := h.Background.Verify(); err != nil {
		return err
	}
	for i := range h.Elements {
		if err := h.Elements[i].Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (d DrawingConfig) Verify() error {
	if err := d.Background.Verify(); err != nil {
		return err
	}
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
