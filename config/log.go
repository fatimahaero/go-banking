package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Global variable to save log file
var logFile *os.File

func InitiateLog() {
	// Makesure folder log exists
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", os.ModePerm)
	}

	// Open or create log file
	var err error
	logFile, err = os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	// Zerolog configuration to write log to file
	log.Logger = zerolog.New(logFile).With().
		Timestamp().
		Logger()

	log.Info().Msg("Zerolog initiated successfully")
}

// Function to close log (called when application is closed)
func CloseLog() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close log file")
		}
	}
}
