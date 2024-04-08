package service

import (
	"log"
	"os"
	"syscall"
)

func Terminate() {
	appProcess, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Printf("failed to safe terminate the app, failed to get the app's pid!!, error: %v", err.Error())
		log.Panic("app is terminating gracefully... ")
	}

	err = appProcess.Signal(syscall.SIGTERM)
	if err != nil {
		log.Printf("failed to safe terminate the app, failed to send the termination signal, error: %v", err.Error())
		log.Panic("app is terminating gracefully... ")
	}

	log.Println("the app is safely terminated")

}
