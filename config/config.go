package config

import (
	"os"
	"sync"
	"time"
)

type Config struct {
	App    App
	Server Server
}

type App struct {
	Name    string
	Version string
}

type Server struct {
	RateLimit      RateLimit
	ReadTimeoutMs  time.Duration
	WriteTimeoutMs time.Duration
	IdleTimeoutMs  time.Duration
}

type RateLimit struct {
	Requests int
	PeriodMs int
}

var (
	once     sync.Once
	Settings Config
)

//This global variable is updated with the provided git tag at release.yml Github Action runtime
var Version = "version"

func InitConfiguration() {
	once.Do(func() {
		//App
		Settings.App.Name = "benthos"
		if os.Getenv("ENV") != "DEV" {
			Settings.App.Version = Version
		} else {
			Settings.App.Version = "v0.0.3"
		}
		//Server
		Settings.Server.ReadTimeoutMs = 15000
		Settings.Server.WriteTimeoutMs = 15000
		Settings.Server.IdleTimeoutMs = 60000
		//Server RateLimit
		Settings.Server.RateLimit.Requests = 10
		Settings.Server.RateLimit.PeriodMs = 10000
	})
}
