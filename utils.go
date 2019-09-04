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

	for k, v := range nameArr {
		nameArr[k] = strings.Title(v)
	}

	return strings.Join(nameArr, "")
}

func getApiNamePrefix(name string) string {
	return strings.Split(name, "_")[0]
}

func makeDefinitionsMap(definitions map[string]responsesDefinition) map[string]map[string]responsesDefinition {
	defMap := make(map[string]map[string]responsesDefinition)

	for k, v := range definitions {
		pref := getApiNamePrefix(k)

		if _, ok := defMap[pref]; ok {
			defMap[pref][k] = v
		} else {
			defMap[pref] = make(map[string]responsesDefinition)
			defMap[pref][k] = v
		}
	}

	return defMap
}

// readSchemaFile: reads VK API schema file from HTTP URL and saves it locally in the working directory
func readSchemaFile(fileUrl string) ([]byte, error) {
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getObjectTypeName(s string) string {
	str := strings.Split(s, "/")
	return strings.Join([]string{"objects", convertName(str[len(str)-1])}, ".")

}
