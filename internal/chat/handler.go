package chat

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"unicode"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"github.com/victormamede/benebott/internal/capabilities"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update, model *genai.GenerativeModel, store *ChatStore) {
	botUser, err := b.GetMe(ctx)
	if err != nil {
		log.Fatal("Could not get user", err)
	}

	if update.Message == nil {
		return
	}

	if isMentionedOrReplied(botUser, update) {
		cs := store.Get(update.Message.Chat.ID, model)

		name := "anonymous"
		if update.Message.From != nil {
			name = update.Message.From.FirstName
		}

		aiCall(ctx, b, update, cs, genai.Text(fmt.Sprintf("[%s] %s", name, update.Message.Text)))

		return
	}

	if rand.Float64() < viper.GetFloat64("bot.mock_chance") && len(update.Message.Text) > 2 {
		mockMessage := []rune{}

		for i, curr := range update.Message.Text {
			if i%2 == 0 {
				mockMessage = append(mockMessage, unicode.ToUpper(curr))
			} else {
				mockMessage = append(mockMessage, unicode.ToLower(curr))
			}
		}

		b.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID:          update.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
				Text:            string(mockMessage),
			})
	}
}

func aiCall(ctx context.Context, b *bot.Bot, update *models.Update, cs *genai.ChatSession, part genai.Part) {
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})

	resp, err := cs.SendMessage(ctx, part)

	if err != nil {
		log.Println("Gemini error", err)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:          update.Message.Chat.ID,
			Text:            "Erro: " + err.Error(),
			ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
		})

		if err != nil {
			log.Println("Reply error", err)
		}
		return
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				switch v := part.(type) {
				case genai.Text:
					_, err := b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:          update.Message.Chat.ID,
						Text:            string(v),
						ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
					})

					if err != nil {
						log.Println("Reply error", err)
						continue
					}

				case genai.FunctionCall:
					response := map[string]any{}

					switch v.Name {
					case capabilities.MyIpDeclaration.Name:
						response = capabilities.MyIp()
					case capabilities.DotaPlayerAccountDeclaration.Name:
						response = capabilities.DotaPlayerAccount(v.Args["playerId"].(string))
					case capabilities.DotaPlayerMatchesDeclaration.Name:
						response = capabilities.DotaPlayerMatches(v.Args["playerId"].(string), int(v.Args["limit"].(float64)))
					case capabilities.DotaHeroesDeclaration.Name:
						response = capabilities.DotaHeroes()
					case capabilities.UnixTimestampDeclaration.Name:
						response = capabilities.UnixTimestamp(int64(v.Args["timestamp"].(float64)))
					}

					aiCall(ctx, b, update, cs, genai.FunctionResponse{
						Name: v.Name, Response: response,
					})
				default:
					log.Println("Unexpected error: ", err)
					continue
				}
			}
		}
	}

}

func isMentionedOrReplied(user *models.User, update *models.Update) bool {
	if update.Message.ReplyToMessage != nil &&
		update.Message.ReplyToMessage.From != nil &&
		update.Message.ReplyToMessage.From.ID == user.ID {
		return true
	}

	for _, entity := range update.Message.Entities {
		if entity.Type == models.MessageEntityTypeMention {
			mentionText := update.Message.Text[entity.Offset : entity.Offset+entity.Length]

			if mentionText == "@"+user.Username {
				return true
			}
		}
	}

	return false
}
