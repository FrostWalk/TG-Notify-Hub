package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"tgnotifyhub/config"
	"tgnotifyhub/healtcheck"
	"tgnotifyhub/markdown"
	"tgnotifyhub/telegram"
)

func Send(c *gin.Context) {
	topicName := c.Param("slug")

	bodyBytes, err := c.GetRawData()
	if err != nil {
		err = telegram.SendMessageToGeneral(config.Loaded().ChatId, formatError(err))
		if err != nil {
			healtcheck.SignalError(err)
			log.Println(err)
		}

		log.Println(err)
		healtcheck.SignalError(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var body string
	if !strings.Contains(c.GetHeader("Content-Type"), "json") {
		body = markdown.Escape(string(bodyBytes))
	} else {
		body = fmt.Sprintf("```\n%s\n```", string(bodyBytes))
	}

	if topicName == "" {
		err = telegram.SendMessageToGeneral(config.Loaded().ChatId, body)
		if err != nil {
			log.Println(err)
			healtcheck.SignalError(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	found, topicId := config.GetIdFromName(topicName)
	if found {
		err = telegram.SendMessageToTopic(config.Loaded().ChatId, topicId, body)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	} else {
		err = telegram.SendMessageToGeneral(config.Loaded().ChatId, wrongName(topicName, body))
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}

func wrongName(slug, text string) string {
	return fmt.Sprintf("⚠️ *Attenzione*\nil topic *%s* non esiste\n\n%s", slug, text)
}

func formatError(error error) string {
	return fmt.Sprintf("*Si è verificato un errore*\n\n%s", error)
}
