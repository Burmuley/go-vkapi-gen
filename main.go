package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	VK_SCHEMA_FILES = map[string]string{
		"VK_API_SCHEMA_OBJECTS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json",
		"VK_API_SCHEMA_METHODS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json",
		"VK_API_SCHEMA_RESPONSES": "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json",
		"RESPONSES_LOCAL":         "/Users/burmuley/go/src/github.com/vk-api-schema/test.json",
	}
)

// readEnvVariables: Read environment variables to override defaults
func readEnvVariables() {
	for k := range VK_SCHEMA_FILES {
		if tmpvar := os.Getenv(k); tmpvar != "" {
			VK_SCHEMA_FILES[k] = tmpvar
		}
	}
}

func main() {
	readEnvVariables()
	fmt.Println(VK_SCHEMA_FILES)

	responses, err := readSchemaFile(VK_SCHEMA_FILES["VK_API_SCHEMA_RESPONSES"])
	//responses, err := readLocalSchemaFile(VK_SCHEMA_FILES["RESPONSES_LOCAL"])

	if err != nil {
		fmt.Println("Error:", err)
	}

	jsonResponses := responsesSchema{}
	if err := json.Unmarshal(responses, &jsonResponses); err != nil {
		fmt.Printf("JSON Error:%#v\n", err)
		return
	}

	generateResponses(jsonResponses)
}
