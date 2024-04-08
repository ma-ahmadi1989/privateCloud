package main

import (
	"logReader/internal/channels"
	"logReader/internal/config"
	"logReader/internal/service"
	"sync"
)

func main() {
	var logReaderWaitGroup sync.WaitGroup
	config.LoadLogReaderConfigs()
	channels.LoadGlobalChannels()

	logReaderWaitGroup.Add(1)
	go service.ReadEvents(&logReaderWaitGroup)

	logReaderWaitGroup.Wait()
}
