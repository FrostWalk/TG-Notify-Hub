package main

import (
	"flag"
	"fmt"
	"github.com/etkecc/go-healthchecks/v2"
	"github.com/gin-gonic/gin"
	"log"
	"tgnotifyhub/api"
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
	if err := telegram.InitBot(config.Loaded().Token); err != nil {
		log.Panic(err)
	}

	// crate new topics
	t, _ := telegram.CreateTopics(config.Loaded().Topics, config.Loaded().ChatId)
	// save topics with ids
	if err := config.UpdateTopics(t, *configFile); err != nil {
		log.Panic(err)
	}

	if config.Loaded().HealthCheckUuid != "" {
		client := healthchecks.New(healthchecks.WithCheckUUID(config.Loaded().HealthCheckUuid))
		defer client.Shutdown()

		go client.Auto(time.Duration(config.Loaded().PingInterval) * time.Second)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(api.AuthMiddleware)
	r.POST("/send/:slug", api.Send)
	r.POST("/send", api.Send)

	if err := r.Run(fmt.Sprintf(":%d", config.Loaded().Port)); err != nil {
		log.Panic(err)
	}
}
