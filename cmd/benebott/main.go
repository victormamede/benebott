package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"unicode"

	"benebott/bot/internal/chat"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/benebott")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Init gemini
	ai, err := genai.NewClient(ctx, option.WithAPIKey(viper.GetString("keys.gemini")))
	if err != nil {
		panic(err)
	}
	defer ai.Close()

	// Init bot
	chat_store := chat.CreateChatStore(viper.GetInt("bot.max_history"))
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			handler(ctx, bot, update, ai, chat_store)
		}),
	}
	b, err := bot.New(viper.GetString("keys.telegram"), opts...)
	if err != nil {
		panic(err)
	}

	// Start bot
	fmt.Println("Bot started..")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update, ai *genai.Client, store *chat.ChatStore) {
	botUser, err := b.GetMe(ctx)
	if err != nil {
		log.Fatal("Could not get user", err)
	}

	if update.Message == nil {
		return
	}

	if isMentionedOrReplied(botUser, update) {
		model := ai.GenerativeModel("gemini-2.0-flash")
		model.SystemInstruction = genai.NewUserContent(genai.Text(viper.GetString("bot.prompt")))

		cs := model.StartChat()
		cs.History = store.Get(update.Message.Chat.ID)

		var message *models.Message = nil
		b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: update.Message.Chat.ID, Action: models.ChatActionTyping})

		name := "anonymous"
		if update.Message.From != nil {
			name = update.Message.From.FirstName
		}

		iter := cs.SendMessageStream(ctx, genai.Text(fmt.Sprintf("[%s] %s", name, update.Message.Text)))
		fullText := ""
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:          update.Message.Chat.ID,
					Text:            err.Error(),
					ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
				})
				return
			}

			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						switch v := part.(type) {
						case genai.Text:
							fullText += string(v)
						default:
							fullText += "[invalid data]"
						}
					}
				}
			}

			if message == nil {
				newMessage, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:          update.Message.Chat.ID,
					Text:            fullText,
					ReplyParameters: &models.ReplyParameters{MessageID: update.Message.ID},
				})

				if err != nil {
					log.Fatal(err)
					continue
				}

				message = newMessage
			} else {
				newMessage, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:    update.Message.Chat.ID,
					MessageID: message.ID,
					Text:      fullText,
				})

				if err != nil {
					log.Fatal(err)
					continue
				}

				message = newMessage
			}
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: message.ID,
			Text:      fullText,
			ParseMode: "MarkdownV2",
		})

		store.Save(update.Message.Chat.ID, cs.History)

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
