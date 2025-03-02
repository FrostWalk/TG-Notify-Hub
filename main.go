package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"tgnotifyhub/api"
	"tgnotifyhub/config"
	"tgnotifyhub/formatters"
	"tgnotifyhub/healtcheck"
	"tgnotifyhub/telegram"
)

const pluginsPath = "./plugins"

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

	if config.Loaded().ChatId == 0 {
		id, err := telegram.GetGroupId()
		if err != nil {
			log.Panic(err)
		}

		err = config.SetGroupId(id)
		if err != nil {
			log.Panic(err)
		}
	}

	// crate new topics
	t, _ := telegram.CreateTopics(config.Loaded().Topics, config.Loaded().ChatId)
	// save topics with ids
	if err := config.UpdateTopics(t); err != nil {
		log.Panic(err)
	}

	// load formatter plugins
	if err := formatters.LoadPluginsFromFolder(pluginsPath); err != nil {
		log.Panic(err)
	}

	if uuid := config.Loaded().HealthCheckUuid; uuid != "" {
		healtcheck.EnableCheck(uuid, config.Loaded().PingInterval)
		defer healtcheck.CloseConnection()
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
