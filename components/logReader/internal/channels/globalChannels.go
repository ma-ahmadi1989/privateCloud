package channels

import (
	"logReader/internal/models"
)

type LogReaderChannels struct {
	GitEvents chan models.GitEvent
}

var GlobalChannels *LogReaderChannels

func init() {
	GlobalChannels = &LogReaderChannels{
		GitEvents: make(chan models.GitEvent, 10),
	}
}

func LoadGlobalChannels() *LogReaderChannels {
	return GlobalChannels
}
