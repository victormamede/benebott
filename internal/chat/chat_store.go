package chat

import (
	"context"

	"github.com/spf13/viper"
	"google.golang.org/genai"
)

type ChatStore struct {
	MaxHistory int

	chats map[int64]*genai.Chat
}

func CreateChatStore(maxHistory int) *ChatStore {
	return &ChatStore{
		chats:      map[int64]*genai.Chat{},
		MaxHistory: maxHistory,
	}
}

func (s *ChatStore) Get(ctx context.Context, id int64, client *genai.Client, config *genai.GenerateContentConfig) *genai.Chat {
	chat, ok := s.chats[id]

	if !ok {
		history := []*genai.Content{}
		chat, _ := client.Chats.Create(ctx, viper.GetString("bot.model"), config, history)
		s.chats[id] = chat

		return chat
	}

	return chat
}
