package capabilities

import "github.com/google/generative-ai-go/genai"

type CallResponse = map[string]any

var Tools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			&MyIpDeclaration,
			&DotaPlayerAccountDeclaration,
			&DotaPlayerMatchesDeclaration,
			&DotaHeroesDeclaration,
			&UnixTimestampDeclaration,
		},
	},
}
