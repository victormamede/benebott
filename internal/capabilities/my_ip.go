package capabilities

import (
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/genai"
)

var MyIpDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "get_my_ip",
	Description: "Gets the ip address of where the bot is currently hosted",
	Parameters: &genai.Schema{
		Type:       genai.TypeObject,
		Properties: map[string]*genai.Schema{},
	},
}

func MyIp() CallResponse {
	response, err := http.Get("https://api.ipify.org?format=json")
	callResponse := map[string]any{}

	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	err = json.Unmarshal(responseData, &callResponse)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	return callResponse
}
