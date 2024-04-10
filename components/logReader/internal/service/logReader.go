package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"logReader/internal/channels"
	"logReader/internal/config"
	"logReader/internal/models"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/tidwall/gjson"
)

func ReadEvents(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigs
		cancel()
		close(channels.GlobalChannels.GitEvents)
		log.Printf("terminate signal recieved on ReadEvent, signal: %+v", sig.String())
	}()

	eventFiles, err := GetEventFilesList()
	if err != nil {
		log.Printf("faile to load event files, error: %s", err.Error())
		log.Println("terminting the app as there is no event to process")
		Terminate()
	}

	for _, eventFile := range eventFiles {
		Extract(eventFile, ctx)
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

func Extract(eventFilePath string, ctx context.Context) {
	eventFile, err := os.Open(eventFilePath)
	if err != nil {
		log.Printf("failed to open event file: %s, error: %v", eventFilePath, err.Error())
		return
	}
	defer func() {
		eventFile.Close()
		recover()
	}()

	fileScanner := bufio.NewScanner(eventFile)

	for {
		select {
		case <-ctx.Done():
			log.Printf("recieving termination signal from os. stop reading event file: %s...", eventFilePath)
			return
		default:
			fileScanner.Scan()
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

}

func ConvertToEventRecord(eventRecord string) (models.GitEvent, error) {
	extractedEvent := models.GitEvent{
		EventID:      gjson.Get(eventRecord, "id").String(),
		Type:         gjson.Get(eventRecord, "type").String(),
		GitUserName:  gjson.Get(eventRecord, "actor.login").String(),
		GitRepoName:  gjson.Get(eventRecord, "repo.name").String(),
		GitRepoURL:   gjson.Get(eventRecord, "repo.url").String(),
		Public:       gjson.Get(eventRecord, "public").Bool(),
		CreatedAt:    gjson.Get(eventRecord, "created_at").Time(),
		CommitsCount: len(gjson.Get(eventRecord, "payload.commits").Array()),
	}

	return extractedEvent, nil
}
