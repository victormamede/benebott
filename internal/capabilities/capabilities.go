package capabilities

import "google.golang.org/genai"

type CallResponse = map[string]any

var Tools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			&MyIpDeclaration,
			&DotaPlayerAccountDeclaration,
			&DotaPlayerMatchesDeclaration,
			&DotaHeroesDeclaration,
			&UnixTimestampDeclaration,
			&MyIdDeclaration,
		},
	},
}
