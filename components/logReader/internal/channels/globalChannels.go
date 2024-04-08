package channels

import (
	"logReader/internal/models"
	"os"
	"os/signal"
	"syscall"
)

type LogReaderChannels struct {
	GitEvents chan models.GitEvent
	OsSignal  chan os.Signal
}

var GlobalChannels *LogReaderChannels

func init() {
	GlobalChannels = &LogReaderChannels{
		GitEvents: make(chan models.GitEvent, 10),
		OsSignal:  make(chan os.Signal, 1),
	}

	signal.Notify(GlobalChannels.OsSignal, syscall.SIGINT, syscall.SIGTERM)
}

func LoadGlobalChannels() *LogReaderChannels {
	return GlobalChannels
}
