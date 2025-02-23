package telegram

import (
	"errors"
	"github.com/mymmrac/telego"
	"net/url"
	"sync"
	"tgnotifyhub/config"
)

var (
	bot *telego.Bot
	mu  sync.Mutex
)

// InitBot initializes the Telegram bot with the provided token.
// It should be called early in your application's startup process.
func InitBot(token string) error {
	var err error
	bot, err = telego.NewBot(token)
	if err != nil {
		return err
	}

	// Log bot initialization details.
	return nil
}

// SendMessageToGeneral sends a text message to the specified chatID.
// This method is safe for concurrent use.
func SendMessageToGeneral(chatID int64, message string) error {
	if bot == nil {
		return errors.New("telegram bot not initialized")
	}

	// Ensure that sending messages is thread-safe.
	mu.Lock()
	defer mu.Unlock()

	// Construct the parameters for sending a message.
	params := &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: chatID},
		Text:   message,
	}

	// Send the message using the telego API.
	_, err := bot.SendMessage(params)
	return err
}

// CreateTopics takes a bot instance, a group chat ID, and an array of topic names.
// For each topic name, it attempts to create a forum topic (if not already created)
// and stores its message thread ID in the topicMap.
func CreateTopics(topics []config.Topic, chatId int64) ([]config.Topic, error) {
	for i, topic := range topics {
		// Skip already created
		if topic.Id != 0 {
			continue
		}

		// Attempt to create the forum topic.
		// Note: The method CreateForumTopic is available in the newer Telegram Bot API.
		// If the topic already exists, the API may return an error; you can decide how to handle that.
		resp, err := bot.CreateForumTopic(&telego.CreateForumTopicParams{
			ChatID: telego.ChatID{ID: chatId},
			Name:   topic.Name,
		})

		if err != nil {
			continue
		}
		// Save the message thread ID associated with the topic.
		topics[i].Id = resp.MessageThreadID
		topics[i].Slug = url.PathEscape(topic.Name)
	}
	return topics, nil
}

// SendMessageToTopic sends a text message to the specified forum topic in the group.
// It looks up the topic name in topicMap to get the message thread ID.
func SendMessageToTopic(chatId int64, topicId int, message string) error {
	// Use the WithMessageThreadID option so the message is sent to the forum topic thread.
	_, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID:          telego.ChatID{ID: chatId},
		Text:            message,
		MessageThreadID: topicId,
	})

	return err
}
