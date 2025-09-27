package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	App      App
	Server   Server
	Database Database
	Security Security
}

type App struct {
	Name    string
	Version string
}

type Server struct {
	RateLimit    RateLimit
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type RateLimit struct {
	Requests int
	Period   time.Duration
}

type Database struct {
	Hostname string
	Port     string
	Username string
	Password string
	Name string
}

type Security struct {
	EncryptionKey string
}

var (
	once     sync.Once
	Settings Config
)

// This global variable is updated with the provided git tag at release.yml Github Action runtime
var Version = "version"

func InitConfiguration() {
	once.Do(func() {
		//App
		Settings.App.Name = "benthos"
		Settings.App.Version = "v0.0.3"
		//Server
		Settings.Server.ReadTimeout = 15000
		Settings.Server.WriteTimeout = 15000
		Settings.Server.IdleTimeout = 60000
		//Server RateLimit
		Settings.Server.RateLimit.Requests = 10
		Settings.Server.RateLimit.Period = 10000
		//DB
		Settings.Database.Hostname = os.Getenv("DB_HOSTNAME")
		Settings.Database.Port = os.Getenv("DB_PORT")
		Settings.Database.Username = os.Getenv("DB_USERNAME")
		Settings.Database.Password = os.Getenv("DB_PASSWORD")
		Settings.Database.Name = os.Getenv("DB_NAME")
		//Security
		Settings.Security.EncryptionKey = os.Getenv("ENCRYPTION_KEY")

		if os.Getenv("ENV") == "QA" || os.Getenv("ENV") == "PROD" {
			//App
			Settings.App.Version = Version
			//DB
			dbPassword, err := loadSecret("db_password")
			if err != nil {
				log.Fatal("DB Password secret not found")
			}
			Settings.Database.Password = dbPassword;
			//Security
			encryptionKey, err := loadSecret("benthos_encryption_key")
			if err != nil {
				log.Fatal("Encryption key secret not found")
			}
			Settings.Database.Password = encryptionKey;
			Settings.Security.EncryptionKey = os.Getenv("ENCRYPTION_KEY")
		}

	})
}

func loadSecret(name string) (string, error) {
	secretPath := filepath.Join("/run/secrets", name)

	if data, err := os.ReadFile(secretPath); err == nil {
		return string(data), nil
	}

	return "", os.ErrNotExist
}
