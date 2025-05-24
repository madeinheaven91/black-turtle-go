package config

import (
	"os"
	"sync"
)

type config struct {
	publicationId string
	botToken      string
	logLevel      string

	pgHost     string
	pgPort     string
	pgDb       string
	pgName     string
	pgUser     string
	pgPassword string
}

var (
	cfg  = config{}
	once sync.Once
)

func InitFromEnv() {
	once.Do(func() {
		cfg.botToken = os.Getenv("BOT_TOKEN")
		if cfg.botToken == "" {
			panic("BOT_TOKEN not provided in .env")
		}
		cfg.publicationId = os.Getenv("PUBLICATION_ID")
		if cfg.publicationId == "" {
			panic("PUBLICATION_ID not provided in .env")
		}
		cfg.logLevel = os.Getenv("LOG_LEVEL")
		if cfg.logLevel == "" {
			cfg.logLevel = "INFO"
		}

		cfg.pgHost = os.Getenv("POSTGRES_HOST")
		if cfg.pgHost == "" {
			panic("POSTGRES_HOST not provided in .env")
		}
		cfg.pgPort = os.Getenv("POSTGRES_PORT")
		if cfg.pgPort == "" {
			panic("POSTGRES_PORT not provided in .env")
		}
		cfg.pgUser = os.Getenv("POSTGRES_USER")
		if cfg.pgUser == "" {
			panic("POSTGRES_USER not provided in .env")
		}
		cfg.pgDb = os.Getenv("POSTGRES_DB")
		if cfg.pgDb == "" {
			panic("POSTGRES_DB not provided in .env")
		}
		cfg.pgPassword = os.Getenv("POSTGRES_PASSWORD")
		if cfg.pgPassword == "" {
			panic("POSTGRES_PASSWORD not provided in .env")
		}
		cfg.pgName = os.Getenv("POSTGRES_NAME")
		if cfg.pgName == "" {
			panic("POSTGRES_NAME not provided in .env")
		}
	})
}

func PublicationID() string {
	return cfg.publicationId
}

func BotToken() string {
	return cfg.botToken
}

func LogLevel() string {
	return cfg.logLevel
}

func PgHost() string {
	return cfg.pgHost
}

func PgPort() string {
	return cfg.pgPort
}

func PgDb() string {
	return cfg.pgDb
}

func PgUser() string {
	return cfg.pgUser
}

func PgName() string {
	return cfg.pgName
}

func PgPassword() string {
	return cfg.pgPassword
}
