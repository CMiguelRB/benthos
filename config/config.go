package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	App    App    `json:"app"`
	Server Server `json:"server"`
}

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Server struct {
	RateLimit RateLimit `json:"rateLimit"`
	ReadTimeoutMs time.Duration `json:"readTimeoutMs"`
	WriteTimeoutMs time.Duration `json:"writeTimeoutMs"`
	IdleTimeoutMs time.Duration `json:"idleTimeoutMs"`
}

type RateLimit struct {
	Requests int `json:"requests"`
	PeriodMs   int `json:"periodMs"`
}

var (
	once     sync.Once
	Settings Config
)

func InitConfiguration() {
	once.Do(func() {
		confFile, err := os.Open("config.json")
		if err != nil {
			log.Fatal("Error opening config.json file!")
		}
		defer confFile.Close()
		confBytes, err := io.ReadAll(confFile)
		if err != nil {
			log.Fatal("Error reading config.json file!")
		}

		err = json.Unmarshal(confBytes, &Settings)
		if err != nil {
			log.Fatal("Error parsing config.json file!")
		}
	})
}
