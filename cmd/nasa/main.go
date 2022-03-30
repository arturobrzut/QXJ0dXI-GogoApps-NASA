package main

import (
	log "github.com/sirupsen/logrus"
	collector "nasa/internal/url_collector"
	"os"
)

func main() {
	log.Info("Start app")
	config, err := collector.Setup(
		readEnvVar("API_KEY", "DEMO_KEY"),
		readEnvVar("PORT", "8080"),
		readEnvVar("CONCURRENT_REQUESTS", "5"))
	if err != nil {
		log.Fatal("Cannot read parameters", err)
	}
	err = collector.Serv(config)
	if err != nil {
		log.Error(err, "Error from Service.")
	}
}

func readEnvVar(envVariable, defaultValue string) string {
	data, found := os.LookupEnv(envVariable)
	if !found {
		return defaultValue
	}
	return data
}
