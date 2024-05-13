package main

import (
	routerv1 "github.com/kitabisa/backend-takehome-test/api/v1"
	"github.com/kitabisa/backend-takehome-test/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	envConfig, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading env config")
	}

	log.Logger = log.With().Caller().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	if envConfig.App.LogPretty == true {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatTimestamp: func(i interface{}) string { return time.Now().Format(time.RFC3339) }})
	}

	switch envConfig.App.LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var dbConn = config.DbConnection{}
	dbConn.InitDbConnectionPool(envConfig.Postgres)

	routerv1.ServeRouter()
}
