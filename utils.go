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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// convertName: function concatenates and titles all words separated by underscore
func convertName(jsonName string) string {
	nameArr := strings.Split(jsonName, "_")

	if nameArr[len(nameArr)-1] == "response" {
		nameArr = nameArr[:len(nameArr)-1]
	}

	// Convert numbers to words according to Golang naming convention
	if strings.Index(nameArr[0], "2") == 0 {
		nameArr[0] = strings.ReplaceAll(nameArr[0], "2", "two")
	}

	for k, v := range nameArr {
		nameArr[k] = strings.Title(v)
	}

	return strings.Join(nameArr, "")
}

func convertParam(param string) string {
	nameArr := strings.Split(param, "_")

	if nameArr[0] == "type" {
		nameArr[0] = "pType"
	}

	for k, v := range nameArr {
		if k != 0 {
			nameArr[k] = strings.Title(v)
		}
	}

	return strings.Join(nameArr, "")
}

func getApiNamePrefix(name string) string {
	return strings.Split(name, "_")[0]
}

func getApiMethodNamePrefix(name string) string {
	return strings.Split(name, ".")[0]
}

func getApiMethodNameSuffix(name string) string {
	return strings.Split(name, ".")[1]
}

// readHTTPSchemaFile: reads VK API schema file from HTTP URL and saves it locally in the working directory
func readHTTPSchemaFile(fileUrl string) ([]byte, error) {
	var schemaFile []byte

	httpResp, err := http.Get(fileUrl)
	defer httpResp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not download from URL %s. Error: %s", fileUrl, err)
	}

	schemaFile, err = ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return nil, fmt.Errorf("could not download from URL %s. Error: %s", fileUrl, err)
	}

	return schemaFile, nil
}

func readLocalSchemaFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func loadSchemaFile(path string) ([]byte, error) {
	if path[:4] != "http" {
		return readLocalSchemaFile(path)
	}

	return readHTTPSchemaFile(path)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getObjectTypeName(s string) string {
	var prefix string

	p := strings.Split(s, "#")

	if len(p[0]) > 0 {
		prefix = strings.Split(p[0], ".")[0]
	}

	str := strings.Split(s, "/")

	if len(prefix) == 0 {
		return convertName(str[len(str)-1])
	}

	return strings.Join([]string{prefix, convertName(str[len(str)-1])}, ".")
}

func logString(s string) {
	log.Println(s)
}

func logJSONError(err error) {
	logString(fmt.Sprintf("JSON Error:%#v\n", err))
}

//func copyStatic(outputDir string) error {
//
//}
