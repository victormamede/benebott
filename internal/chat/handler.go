package chat

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"unicode"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/viper"
	"github.com/victormamede/benebott/internal/capabilities"
	"google.golang.org/genai"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update, aiClient *genai.Client, config *genai.GenerateContentConfig, store *ChatStore) {
	botUser, err := b.GetMe(ctx)
	if err != nil {
		log.Fatal("Could not get user", err)
	}

	if update.Message == nil {
		return
	}

	if isMentionedOrReplied(botUser, update) {
		cs := store.Get(ctx, update.Message.Chat.ID, aiClient, config)

		name := "anonymous"
		if update.Message.From != nil {
			name = update.Message.From.FirstName
		}

		aiCall(ctx, b, update, cs, *genai.NewPartFromText(fmt.Sprintf("[%s] %s", name, update.Message.Text)))

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

		return
	}

}

func aiCall(ctx context.Context, b *bot.Bot, update *models.Update, chat *genai.Chat, part genai.Part) {
	b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})

	resp, err := chat.SendMessage(ctx, part)

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

				if part.Text != "" {
					_, err := b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:          update.Message.Chat.ID,
						Text:            part.Text,
						ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
					})

					if err != nil {
						log.Println("Reply error", err)
						continue
					}

				} else if part.FunctionCall != nil {
					response := map[string]any{}
					v := part.FunctionCall

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

					aiCall(ctx, b, update, chat, genai.Part{
						FunctionResponse: &genai.FunctionResponse{
							Name: v.Name, Response: response,
						},
					})
				} else {
					log.Println("Unexpected error: ", err)
					continue

				}
			}
		}
	}

}

func translateUnintelligible(ctx context.Context, b *bot.Bot, update *models.Update, aiClient *genai.Client, config *genai.GenerateContentConfig) {

	history := []*genai.Content{
		genai.NewContentFromText(viper.GetString("bot.unintelligible_prompt"), genai.RoleUser),
		genai.NewContentFromText(update.Message.Text, genai.RoleUser),
		genai.NewContentFromText("Great to meet you. What would you like to know?", genai.RoleModel),
	}

	resp, err := aiClient.Models.GenerateContent(
		ctx,
		viper.GetString("bot.model"),
		history,
		config,
	)

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
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:          update.Message.Chat.ID,
					Text:            part.Text,
					ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
				})

				if err != nil {
					log.Println("Reply error", err)
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
