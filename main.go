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
	msg   string // Step message to print in log
	fName string // key name to find schema file name
	sObj  IGenerator // IGenerator object to process
}

var (
    // vkSchemaFiles - map of strings to define schema files paths
	vkSchemaFiles = map[string]string{
		"VK_API_SCHEMA_OBJECTS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json",
		"VK_API_SCHEMA_METHODS":   "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json",
		"VK_API_SCHEMA_RESPONSES": "https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json",
	}

	// vkSteps - list of steps to perform to generate resulting VK SDK code
	vkSteps = []step{
		// responses depends on objects
		{"Generating VK API objects", "VK_API_SCHEMA_OBJECTS", &objectsSchema{}},
		{"Generating VK API responses", "VK_API_SCHEMA_RESPONSES", &responsesSchema{}},
		{"Generating VK API methods", "VK_API_SCHEMA_METHODS", &schemaMethods{}},
	}

	// copy objects container to render `allOf` and `oneOf` properties in responses
	objectsGlobal *objectsSchema

	// hash of output directories names
	outputDirs = []string{
		fmt.Sprintf("%s/objects", outputDirName),
		fmt.Sprintf("%s/responses", outputDirName),
	}
)

// readEnvVariables: Read environment variables and override defaults if found
func readEnvVariables() {
	for k := range vkSchemaFiles {
		if tmp := os.Getenv(k); tmp != "" {
			vkSchemaFiles[k] = tmp
		}
	}

	if tmp := os.Getenv("VK_API_SCHEMA_OUTPUT"); tmp != "" {
		outputDirName = tmp
	}
}

// printEnvInfo: print runtime environment information
func printEnvInfo() {
	logInfo("Running with the following configuration parameters:")

	for k, v := range vkSchemaFiles {
		logInfo(fmt.Sprintf("%s = %s", k, v))
	}
}

func main() {
	readEnvVariables()
	printEnvInfo()

	// check and create output directories
	if err := makeDirs(outputDirs); err != nil {
		logError(err)
		os.Exit(1)
	}

	// copy static code to the output directory
	if err := copyStatic(outputDirName); err != nil {
		logError(err)
		os.Exit(1)
	} else {
		logInfo("static content copied successfully")
	}

	// walk through steps and run parse&generate methods
	for _, v := range vkSteps {
		logStep(v.msg)

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
