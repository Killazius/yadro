package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	configPath = "sunny_5_skiers/config.json"
)

// Config - структура для хранения конфигурации гонки по биатлону
type Config struct {
	Laps        int    `json:"laps"`        // Количество кругов в гонке
	LapLen      int    `json:"lapLen"`      // Длина одного круга
	PenaltyLen  int    `json:"penaltyLen"`  // Длина штрафного круга
	FiringLines int    `json:"firingLines"` // Количество стрельбищ
	Start       string `json:"start"`       // Время начала гонки
	StartDelta  string `json:"startDelta"`  // Интервал между стартами участников
}

func MustLoad() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	return &cfg
}
