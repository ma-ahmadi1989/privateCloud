package service

import (
	"bytes"
	"context"
	"encoding/json"
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

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	var eventLoadWaitGroup sync.WaitGroup
	delayBetweenRequests := 1000000 / config.LogReaderConfigs.Rate // micsosecond

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
			go Load(event, ctx, &eventLoadWaitGroup)
			time.Sleep(time.Duration(delayBetweenRequests) * time.Microsecond)
		default:
			log.Println("no event to load")
			time.Sleep(1 * time.Second)
		}
	}

}

func Load(eventLog models.GitEvent, ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		recover()
		wg.Done()
	}()

	jsonEvent, _ := json.Marshal(eventLog)
	req, err := http.NewRequest("POST", config.LoadLogReaderConfigs().GatewayURL, bytes.NewBuffer(jsonEvent))
	if err != nil {
		log.Fatalf("Error occurred during request creation: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occurred during request execution: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error occurred while reading the response body: %s", err.Error())
	}

	log.Printf("Response: %s", body)

}
