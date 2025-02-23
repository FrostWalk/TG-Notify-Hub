package main

import (
	"flag"
	"github.com/etkecc/go-healthchecks/v2"
	"log"
	"tgnotifyhub/config"
	"tgnotifyhub/telegram"
	"time"
)

func main() {
	configFile := flag.String("s", "settings.json", "Path to settings.json")
	flag.Parse()

	// load config from file
	if err := config.Load(*configFile); err != nil {
		log.Panic(err)
	}

	// initialize bot instance
	if err := telegram.InitBot(config.Config().Token); err != nil {
		log.Panic(err)
	}

	// crate new topics
	t, _ := telegram.CreateTopics(config.Config().Topics, config.Config().ChatId)
	// save topics with ids
	if err := config.UpdateTopics(t, *configFile); err != nil {
		log.Panic(err)
	}

	if config.Config().HealthCheckUuid != "" {
		client := healthchecks.New(healthchecks.WithCheckUUID(config.Config().HealthCheckUuid))
		defer client.Shutdown()

		go client.Auto(time.Duration(config.Config().PingInterval) * time.Second)
	}
}
