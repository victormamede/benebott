# Benebott 🤖✨

**Benebott** is a Telegram bot written in Go that integrates with the [Gemini](https://deepmind.google/technologies/gemini/) API. It's built with modularity and configurability in mind, using [Viper](https://github.com/spf13/viper) for flexible configuration management.

## 🚀 Features

- Telegram Bot integration using Go
- Connects to Gemini API for powerful AI functionality
- Modular code structure (`cmd`, `internal`)
- Configurable via TOML, YAML, JSON, etc. (powered by Viper)

## 🛠️ Building the Project

To build the bot, run:

```bash
go build -o ./bin ./cmd/benebott
```

This will compile the binary into the `./bin` directory.

## 📁 Project Structure

```bash
.
├── cmd/
│   └── benebott/
│       └── main.go       # Main entrypoint of the application
├── configs/
│   └── config.toml       # Sample configuration file
├── internal/             # Core libraries and internal logic
└── README.md             # You're here!
```

## ⚙️ Configuration

Benebott uses [Viper](https://github.com/spf13/viper) for configuration, supporting multiple formats (TOML, YAML, JSON, etc).

### 🔍 Configuration file search path

1. Current working directory
2. `/etc/benebott/`

### 📄 Sample config

A sample configuration file is provided at:

```bash
configs/config.toml
```

You can either copy it to the working directory or move it to `/etc/benebott/config.toml`.

## 🧪 Running the Bot

Once built and configured, simply run the bot binary:

```bash
./bin/benebott
```

Make sure your configuration file is properly set up and placed in one of the supported paths.

## 🧱 Dependencies

- [Go](https://golang.org/)
- [Viper](https://github.com/spf13/viper)
- [Golang Telegram Bot](https://github.com/go-telegram/bot)

## 📜 License

MIT License. See [LICENSE](./LICENSE) for details.
