package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Laps        int       `json:"laps"`
	LapLen      float64   `json:"lapLen"`
	PenaltyLen  float64   `json:"penaltyLen"`
	FiringLines int       `json:"firingLines"`
	Start       string    `json:"start"`
	StartDelta  string    `json:"startDelta"`
	StartTime   time.Time // parsed
	Delta       time.Duration
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.StartTime, err = time.Parse("15:04:05.000", cfg.Start)
	if err != nil {
		return nil, err
	}
	deltaParsed, err := time.Parse("15:04:05", cfg.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse StartDelta: %w", err)
	}
	cfg.Delta = time.Duration(deltaParsed.Hour())*time.Hour +
		time.Duration(deltaParsed.Minute())*time.Minute +
		time.Duration(deltaParsed.Second())*time.Second
	return &cfg, nil
}
