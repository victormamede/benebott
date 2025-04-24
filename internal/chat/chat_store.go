package chat

import (
	"github.com/google/generative-ai-go/genai"
)

type ChatSession []*genai.Content

type ChatStore struct {
	MaxHistory int

	sessions map[int64]ChatSession
}

func CreateChatStore(maxHistory int) *ChatStore {
	return &ChatStore{
		sessions:   map[int64]ChatSession{},
		MaxHistory: maxHistory,
	}
}

func (s *ChatStore) Get(id int64) ChatSession {
	session, ok := s.sessions[id]

	if !ok {
		return []*genai.Content{}
	}

	return session
}

func (s *ChatStore) Save(id int64, history []*genai.Content) {
	for len(history) > s.MaxHistory {
		history = history[1:]
	}

	s.sessions[id] = history
}
