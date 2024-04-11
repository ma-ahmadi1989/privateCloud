package config

import (
	"fmt"
	"os"
)

type API struct {
	ListenAddress string
	ListenPort    string
}

func LoadAPIConfig() API {
	api := API{}

	api.ListenAddress = os.Getenv("API_LISTEN_ADDRESS")
	if api.ListenAddress == "" {
		api.ListenAddress = "127.0.0.1"
	}

	api.ListenPort = os.Getenv("API_LISTEN_PORT")
	if api.ListenPort == "" {
		api.ListenPort = "9898"
	}

	return api
}

func (api *API) URI() string {
	return fmt.Sprintf("%s:%s", api.ListenAddress, api.ListenPort)
}
