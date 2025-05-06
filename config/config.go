package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
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

// MustLoad - проверяет есть ли файл по пути path и записывает данные в структуру Config
func MustLoad(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("error loading config")
	}
	return &cfg
}
