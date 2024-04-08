package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"logReader/internal/channels"
	"logReader/internal/config"
	"logReader/internal/models"
	"os"
	"path/filepath"
	"sync"
)

func ReadEvents(wg *sync.WaitGroup) {
	defer func() {
		close(channels.GlobalChannels.GitEvents)
		wg.Done()
	}()

	eventFiles, err := GetEventFilesList()
	if err != nil {
		log.Printf("faile to load event files, error: %s", err.Error())
		log.Println("terminting the app as there is no event to process")
		Terminate()
	}

	for _, eventFile := range eventFiles {
		Extract(eventFile)
	}

}

func GetEventFilesList() ([]string, error) {
	var eventFilesList []string
	eventFiles, err := os.ReadDir(config.LogReaderConfigs.EventFilesPath)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to list the event files, error: %v", err.Error())
		return nil, errors.New(errorMessage)
	}

	for _, eventFileName := range eventFiles {
		if filepath.Ext(eventFileName.Name()) == ".json" {
			eventFilesList = append(eventFilesList,
				filepath.Join(config.LogReaderConfigs.EventFilesPath, eventFileName.Name()))
		}
	}

	return eventFilesList, nil
}

func Extract(eventFilePath string) {
	eventFile, err := os.Open(eventFilePath)
	if err != nil {
		log.Printf("failed to open event file: %s, error: %v", eventFilePath, err.Error())
		return
	}

	fileScanner := bufio.NewScanner(eventFile)

	for fileScanner.Scan() {
		if fileScanner.Err() != nil {
			log.Printf("failed to convert event, error: %v", fileScanner.Err().Error())
			continue
		}

		eventRecord, err := ConvertToEventRecord(fileScanner.Text())
		if err != nil {
			log.Printf("failed to convert event, error: %v", err.Error())
			continue
		}
		channels.GlobalChannels.GitEvents <- eventRecord
	}
}

func ConvertToEventRecord(eventRecord string) (models.GitEvent, error) {
	var extractedEvent models.GitEvent
	var loadedEvent map[string]interface{}
	err := json.Unmarshal([]byte(eventRecord), loadedEvent)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to conver the event record: %s, error: %v", eventRecord, err.Error())
		log.Println(errorMessage)
		return extractedEvent, errors.New(errorMessage)
	}

	extractedEvent = models.GitEvent{
		EventID:     loadedEvent["id"].(string),
		Type:        loadedEvent["type"].(string),
		GitUserName: loadedEvent["actor"]["login"].string,
	}

}
