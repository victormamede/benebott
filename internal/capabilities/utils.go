package capabilities

import (
	"time"

	"github.com/google/generative-ai-go/genai"
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
