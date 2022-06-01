package s4

import (
	"encoding/json"
	"os"
)

type Config struct {
	Input     string `json:"input"`
	Output    string `json:"output"`
	Message   string `json:"message"`
	Key       string `json:"key"`
	Seed      int    `json:"seed"`
	Selection int    `json:"selection"`
}

func LoadConfiguration(file string) (*Config, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	err = json.NewDecoder(configFile).Decode(config)
	return config, err
}
