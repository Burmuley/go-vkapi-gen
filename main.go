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
	"fmt"
	"os"
)

type step struct {
	fName string
	sObj  IGenerator
}

var (
	vkSchemaFiles = map[string]string{
		"VK_API_SCHEMA_OBJECTS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json",
		"VK_API_SCHEMA_METHODS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json",
		"VK_API_SCHEMA_RESPONSES": "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json",
	}

	vkSteps = map[string]step{
		"Generating VK API objects":   step{"VK_API_SCHEMA_OBJECTS", &objectsSchema{}},
		"Generating VK API methods":   step{"VK_API_SCHEMA_METHODS", &schemaMethods{}},
		"Generating VK API responses": step{"VK_API_SCHEMA_RESPONSES", &responsesSchema{}},
	}
)

// readEnvVariables: Read environment variables to override defaults
func readEnvVariables() {
	for k := range vkSchemaFiles {
		if tmp := os.Getenv(k); tmp != "" {
			vkSchemaFiles[k] = tmp
		}
	}
}

func printEnvInfo() {
	logInfo("Running with the following configuration parameters:")

	for k, v := range vkSchemaFiles {
		logInfo(fmt.Sprintf("%s = %s", k, v))
	}
}

func main() {
	readEnvVariables()
	printEnvInfo()

	if err := copyStatic(outputDirName); err != nil {
		logError(err)
		os.Exit(1)
	} else {
		logInfo("static content copied successfully")
	}

	for k, v := range vkSteps {
		logStep(k)

		if err := v.sObj.Parse(vkSchemaFiles[v.fName]); err != nil {
			logError(err)
			os.Exit(1)
		}

		if err := v.sObj.Generate(outputDirName); err != nil {
			logError(err)
			os.Exit(1)
		}
	}
}
