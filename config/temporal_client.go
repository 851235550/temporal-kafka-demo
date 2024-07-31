package config

import (
	"go.temporal.io/sdk/client"
)

func NewTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
}
