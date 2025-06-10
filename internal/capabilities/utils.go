package capabilities

import (
	"time"

	"github.com/go-telegram/bot/models"
	"google.golang.org/genai"
)

var UnixTimestampDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "unix_timestamp",
	Description: "converts an integer to date and time using unix timestamps",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"timestamp": &genai.Schema{Type: genai.TypeInteger},
		},
		Required: []string{"timestamp"},
	},
}

func UnixTimestamp(timestamp int64) CallResponse {
	tm := time.Unix(timestamp, 0)

	result := CallResponse{"time": tm.Format(time.RFC3339)}

	return result
}

var MyIdDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "my_id",
	Description: "returns the id of the user that sent the message",
	Parameters: &genai.Schema{
		Type:       genai.TypeObject,
		Properties: map[string]*genai.Schema{},
	},
}

func MyId(update *models.Update) CallResponse {
	return CallResponse{"id": update.Message.From.ID}
}
