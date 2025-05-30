package chat

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
)

type ChatStore struct {
	MaxHistory int

	sessions map[int64]*genai.ChatSession
}

func CreateChatStore(maxHistory int) *ChatStore {
	return &ChatStore{
		sessions:   map[int64]*genai.ChatSession{},
		MaxHistory: maxHistory,
	}
}

func (s *ChatStore) Get(id int64, model *genai.GenerativeModel) *genai.ChatSession {
	session, ok := s.sessions[id]

	if !ok {
		session = model.StartChat()
		s.sessions[id] = session

		return session
	}

	historyLength := len(session.History)
	maxHistory := viper.GetInt("bot.max_history")
	if historyLength > 10 {
		session.History = session.History[historyLength-maxHistory:]
	}

	return session
}
