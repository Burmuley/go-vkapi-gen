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
		"RESPONSES_LOCAL":         "/Users/konstantin_vasilev/go/src/github.com/vk-api-schema/resp_test_min.json",
		"OBJECTS_LOCAL":           "/Users/konstantin_vasilev/go/src/github.com/vk-api-schema/obj_test_min.json",
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

func printEnvInfo() {
	fmt.Println("Running with the following configuration parameters:")

	for k, v := range VK_SCHEMA_FILES {
		fmt.Printf("%s = %s\n", k, v)
	}
}

func main() {
	readEnvVariables()
	printEnvInfo()

	//responses, err := loadSchemaFile(VK_SCHEMA_FILES["RESPONSES_LOCAL"])

	//responses, err := loadSchemaFile(VK_SCHEMA_FILES["VK_API_SCHEMA_RESPONSES"])
	//
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	//objects, err := loadSchemaFile(VK_SCHEMA_FILES["VK_API_SCHEMA_OBJECTS"])
	objects, err := loadSchemaFile(VK_SCHEMA_FILES["OBJECTS_LOCAL"])

	if err != nil {
		fmt.Println("Error:", err)
	}

	//jsonResponses := responsesSchema{}
	//if err := json.Unmarshal(responses, &jsonResponses); err != nil {
	//	logJSONError(err)
	//	return
	//}

	//for k, v := range jsonResponses.Definitions {
	//	fmt.Println("\n", k)
	//	if r, ok := v.Properties["response"]; ok {
	//		fmt.Println("\t", r.GetGoType(), r.GetType())
	//		if r.GetType() == SCHEMA_TYPE_OBJECT {
	//			for k, v := range r.GetProperties() {
	//				fmt.Println("\t\t", k, v)
	//			}
	//		}
	//	}
	//
	//}

	//generateResponses(jsonResponses)

	jsonObjects := objectsSchema{}

	if err := json.Unmarshal(objects, &jsonObjects); err != nil {
		fmt.Printf("JSON Error:%#v\n", err)
		return
	}

	generateObjects(jsonObjects)
}
