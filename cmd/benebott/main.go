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
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
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
	model := ai.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text(viper.GetString("bot.prompt")))
	model.Tools = capabilities.Tools

	chat_store := chat.CreateChatStore(viper.GetInt("bot.max_history"))
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			chat.Handler(ctx, bot, update, model, chat_store)
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
