/*
Copyright 2019 Konstantin Vasilev (burmuley@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	//
	//objects, err := loadSchemaFile(VK_SCHEMA_FILES["VK_API_SCHEMA_OBJECTS"])
	////objects, err := loadSchemaFile(VK_SCHEMA_FILES["OBJECTS_LOCAL"])
	//
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}
	//
	//jsonResponses := responsesSchema{}
	//if err := json.Unmarshal(responses, &jsonResponses); err != nil {
	//	logJSONError(err)
	//	return
	//}
	//
	//generateResponses(jsonResponses)

	//jsonObjects := objectsSchema{}
	//
	//if err := json.Unmarshal(objects, &jsonObjects); err != nil {
	//	fmt.Printf("JSON Error:%s\n", err)
	//	return
	//}
	//
	//generateObjects(jsonObjects)

	methods, _ := loadSchemaFile(VK_SCHEMA_FILES["VK_API_SCHEMA_METHODS"])

	jsonMethods := schemaMethods{}

	if err := json.Unmarshal(methods, &jsonMethods); err != nil {
		logJSONError(err)
		return
	}

	imethods := make([]IMethod, 0)

	for _, v := range jsonMethods.Methods {
		imethods = append(imethods, v)
	}

	//for _, v := range jsonMethods.Methods {
	//	if v.IsExtended() {
	//		fmt.Println(v.GetName())
	//	}
	//}

	generateMethods(imethods)
}
