package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type LogReaderConfig struct {
	EventFilesPath string `example:"/var/git_events/"`      // This field determines the path where the event files are stored
	Rate           int    `example:"5"`                     // This field determines the rate of the events per second that is supposed to be sent to the service gateway
	GatewayURL     string `example:"http://127.0.0.1:9898"` // This field determines the URL where the requests are sent to
}

var LogReaderConfigs *LogReaderConfig

func init() {
	var err error
	LogReaderConfigs = &LogReaderConfig{}
	if err := godotenv.Load("config.env"); err != nil {
		log.Printf("failed to load the config file: %s, error: %s, default configs will be used", "config.env", err.Error())
	}

	LogReaderConfigs.EventFilesPath = os.Getenv("EVENT_FILES_PATH")
	if LogReaderConfigs.EventFilesPath == "" {
		LogReaderConfigs.EventFilesPath = "events/"
	}
	if !strings.HasPrefix(LogReaderConfigs.EventFilesPath, "/") {
		LogReaderConfigs.EventFilesPath += "/"
	}

	LogReaderConfigs.Rate, err = strconv.Atoi(os.Getenv("EVENT_UPLOAD_RATE"))
	if err != nil {
		LogReaderConfigs.Rate = 5
	}

	LogReaderConfigs.GatewayURL = os.Getenv("GIT_INSIGHT_GATEWAY_URl")
	if LogReaderConfigs.GatewayURL == "" {
		LogReaderConfigs.GatewayURL = "http://127.0.0.1:9898"
	}
}

func LoadLogReaderConfigs() *LogReaderConfig {
	return LogReaderConfigs
}
