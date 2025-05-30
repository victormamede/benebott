package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/victormamede/benebott/internal/capabilities"
	"github.com/victormamede/benebott/internal/chat"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/viper"
	"google.golang.org/genai"
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
	aiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  viper.GetString("keys.gemini"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(viper.GetString("bot.prompt"), genai.RoleUser),
		Tools:             capabilities.Tools,
	}

	chat_store := chat.CreateChatStore(viper.GetInt("bot.max_history"))
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			chat.Handler(ctx, bot, update, aiClient, config, chat_store)
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
