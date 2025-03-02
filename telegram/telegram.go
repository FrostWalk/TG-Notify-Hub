package telegram

import (
	"errors"
	"github.com/mymmrac/telego"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"sync"
	"tgnotifyhub/config"
)

var (
	bot          *telego.Bot
	mu           sync.Mutex
	botInitError = errors.New("telegram bot not initialized")
)

// InitBot initializes the Telegram bot with the provided token.
// It should be called early in your application's startup process.
func InitBot(token string) error {
	var err error
	bot, err = telego.NewBot(token)
	if err != nil {
		return err
	}

	return nil
}

// SendMessageToGeneral sends a text message to the specified chatID.
// This method is safe for concurrent use.
func SendMessageToGeneral(chatID int64, message string) error {
	if bot == nil {
		return botInitError
	}

	// Ensure that sending messages is thread-safe.
	mu.Lock()
	defer mu.Unlock()

	// Construct the parameters for sending a message.
	params := &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		Text:      message,
		ParseMode: "MarkdownV2",
	}

	// Send the message using the telego API.
	_, err := bot.SendMessage(params)
	return err
}

// CreateTopics takes a bot instance, a group chat ID, and an array of topic names.
// For each topic name, it attempts to create a forum topic (if not already created)
// and stores its message thread ID in the topicMap.
func CreateTopics(topics []config.Topic, chatId int64) ([]config.Topic, error) {
	if bot == nil {
		return nil, botInitError
	}

	for i, topic := range topics {
		// Skip already created
		if topic.Id != 0 {
			continue
		}

		// Attempt to create the forum topic.
		// Note: The method CreateForumTopic is available in the newer Telegram Bot API.
		// If the topic already exists, the API may return an error; you can decide how to handle that.
		resp, err := bot.CreateForumTopic(&telego.CreateForumTopicParams{
			ChatID:    telego.ChatID{ID: chatId},
			Name:      topic.Name,
			IconColor: randomColor(),
		})

		if err != nil {
			continue
		}
		// Save the message thread ID associated with the topic.
		topics[i].Id = resp.MessageThreadID
		topics[i].Slug = strings.ToLower(url.PathEscape(topic.Name))
	}
	return topics, nil
}

// SendMessageToTopic sends a text message to the specified forum topic in the group.
// It looks up the topic name in topicMap to get the message thread ID.
func SendMessageToTopic(chatId int64, topicId int, message string) error {
	if bot == nil {
		return botInitError
	}

	// Use the WithMessageThreadID option so the message is sent to the forum topic thread.
	_, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID:          telego.ChatID{ID: chatId},
		Text:            message,
		MessageThreadID: topicId,
		ParseMode:       "MarkdownV2",
	})

	return err
}

func randomColor() int {
	const size = 6
	var colors = [size]int{7322096, 16766590, 13338331, 9367192, 16749490, 16478047}
	return colors[rand.Intn(size)]
}

func GetGroupId() (int64, error) {
	if bot == nil {
		return 0, botInitError
	}

	updates, err := bot.GetUpdates(&telego.GetUpdatesParams{
		Offset: 0,
	})
	if err != nil {
		return 0, err
	}

	sort.Slice(updates, func(i, j int) bool {
		return updates[i].UpdateID < updates[j].UpdateID
	})

	for i := len(updates) - 1; i > 0; i-- {
		u := updates[i]
		if u.MyChatMember != nil && u.MyChatMember.Chat.IsForum {
			return u.MyChatMember.Chat.ID, nil
		}
	}

	return 0, errors.New("unable to determine group id")
}
