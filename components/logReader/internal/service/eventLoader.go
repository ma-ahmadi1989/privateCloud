package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"logReader/internal/channels"
	"logReader/internal/config"
	"logReader/internal/models"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func EventLoader(wg *sync.WaitGroup) {
	defer func() {
		if recoverEvent := recover(); recoverEvent != nil {
			log.Println(recoverEvent)
		}
		wg.Done()
	}()

	var eventLoadWaitGroup sync.WaitGroup
	delayBetweenRequests := 1000000 / config.LogReaderConfigs.Rate // microseconds
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigs
		log.Printf("terminate signal recieved on  EventLoader, signal: %+v", sig.String())
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-channels.GlobalChannels.GitEvents:
			eventLoadWaitGroup.Add(1)
			go Load(event, &eventLoadWaitGroup)
			time.Sleep(time.Duration(delayBetweenRequests) * time.Microsecond)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func Load(eventLog models.GitEvent, wg *sync.WaitGroup) {
	defer func() {
		recover()
		wg.Done()
	}()

	req, err := GenerateRequest(eventLog)
	if err != nil {
		log.Printf("Error occurred during request creation: %s", err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error occurred during request execution: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error occurred while reading the response body: %s", err.Error())
	}

	log.Printf("Response: %s, status: %s (%d)", body, resp.Status, resp.StatusCode)

}

func GenerateRequest(eventLog models.GitEvent) (*http.Request, error) {
	evenLogJson, _ := json.Marshal(eventLog)
	httpRequest, err := http.NewRequest("POST", config.LogReaderConfigs.GatewayURL, bytes.NewBuffer(evenLogJson))
	if err != nil {
		errorMessage := fmt.Sprintf("failed to marshal event log before loading, error: %+v", err.Error())
		return httpRequest, errors.New(errorMessage)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	return httpRequest, nil
}
