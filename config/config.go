package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type AppConfig struct {
	Token           string  `json:"token"`
	ChatId          int64   `json:"chat_id"`
	Port            int     `json:"port"`
	HealthCheckUuid string  `json:"healthcheck_uuid"`
	PingInterval    int     `json:"ping_interval"`
	AuthHeader      string  `json:"auth_header"`
	AuthToken       string  `json:"auth_token"`
	Topics          []Topic `json:"topics"`
}
type Topic struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
	Id   int    `json:"id,omitempty"`
}

var (
	// instance holds the application configuration.
	instance *AppConfig
	// mu protects access to instance.
	mu sync.RWMutex
	// slugChatIds contains the association between slugs and chat ids
	slugChatIds sync.Map
)

// Load attempts to read the configuration from the provided file path.
// If the file does not exist, it creates a default configuration file, logs the event,
// and exits the application.
func Load(path string) error {
	// Check if the file exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create a default configuration.
		defaultConfig := AppConfig{
			Port: 8080,
			Topics: []Topic{
				{
					Name: "",
				},
			},
		}

		// Marshal the default configuration to JSON.
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return err
		}

		// Write the default configuration to file.
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}

		// Log that the config file was created and exit the application.
		log.Printf("Configuration file not found. A default config file has been created at %s. Please update it as needed and restart the application.", path)
		os.Exit(0)
	}

	// File exists; read its content.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into our AppConfig struct.
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Store the configuration globally.
	mu.Lock()
	instance = &cfg
	mu.Unlock()

	for _, topic := range cfg.Topics {
		slugChatIds.Store(topic.Slug, topic.Id)
	}

	return nil
}

// saveConfig writes the current configuration into a JSON file.
func saveConfig(path string) error {
	mu.RLock()
	data, err := json.MarshalIndent(instance, "", "  ")
	mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Config returns a copy of the current configuration.
// All parts of the application can use this to access the fields.
func Config() *AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func UpdateTopics(t []Topic, f string) error {
	instance.Topics = t

	for _, topic := range t {
		slugChatIds.Store(topic.Slug, topic.Id)
	}

	return saveConfig(f)
}
