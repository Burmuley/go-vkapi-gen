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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// convertName: function concatenates and titles all words separated by underscore
func convertName(jsonName string) string {
	nameArr := strings.Split(jsonName, "_")

	// Convert numbers to words according to Golang naming convention
	if strings.Index(nameArr[0], "2") == 0 {
		nameArr[0] = strings.ReplaceAll(nameArr[0], "2", "two")
	}

	for k, v := range nameArr {
		nameArr[k] = strings.Title(strings.ToLower(v))
	}

	return strings.Join(nameArr, "")
}

func cutSuffix(str, suf string) string {
	// don't cut "Response" suffix if it's from objects package
	if strings.Count(str, "objects.") > 0 && suf == "Response" {
		return str
	}

	return strings.TrimSuffix(str, suf)
}

//func cutPrefix(str, pref string) string {
//	return strings.TrimPrefix(str, pref)
//}

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
	var sep string

	if strings.Count(name, ".") > 0 {
		sep = "."
	} else if strings.Count(name, "_") > 0 {
		sep = "_"
	}

	return strings.Split(name, sep)[0]
}

//func getApiMethodNamePrefix(name string) string {
//	return strings.Split(name, ".")[0]
//}

func getApiMethodNameSuffix(name string) string {
	return strings.Split(name, ".")[1]
}

// readHTTPSchemaFile: reads VK API schema file from HTTP URL and saves it locally in the working directory
func readHTTPSchemaFile(fileUrl string) ([]byte, error) {
	logInfo(fmt.Sprintf("Downloading schema file from '%s'", fileUrl))
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
	logInfo(fmt.Sprintf("Loading schema file from '%s'", filePath))
	return ioutil.ReadFile(filePath)
}

func loadSchemaFile(path string) ([]byte, error) {
	if path[:4] != "http" {
		return readLocalSchemaFile(path)
	}

	return readHTTPSchemaFile(path)
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

// Logging helpers
func logString(s string) {
	log.Println(s)
}

func logError(err error) {
	logString(fmt.Sprintf("[ERROR] %#v\n", err))
}

func logInfo(s string) {
	logString(fmt.Sprintf("[INFO] %s", s))
}

func logStep(s string) {
	logInfo(fmt.Sprintf("STEP - %s", s))
}

func checkFileExists(f string) bool {
	fInfo, _ := os.Stat(f)
	return fInfo != nil
}

// Files operations
func copyStatic(outputDir string) error {
	logStep(fmt.Sprintf("Copying static content from `static` directory to `%s`", outputDir))
	staticDir := "./static/"
	return copyDir(staticDir, outputDir)
}

func copyDir(src string, dst string) error {
	var (
		err      error
		fileObjs []os.FileInfo
		srcInfo  os.FileInfo
	)

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fileObjs, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fileObjs {
		srcFileObj := path.Join(src, fd.Name())
		dstFileObj := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcFileObj, dstFileObj); err != nil {
				logError(err)
			}
		} else {
			if err = copyFile(srcFileObj, dstFileObj); err != nil {
				logError(err)
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	var (
		err     error
		srcFile *os.File
		dstFile *os.File
		srcInfo os.FileInfo
	)

	defer srcFile.Close()
	defer dstFile.Close()

	if srcFile, err = os.Open(src); err != nil {
		return err
	}

	if dstFile, err = os.Create(dst); err != nil {
		return err
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

func detectGoType(s string) string {
	switch s {
	case schemaTypeNumber:
		return "json.Number"
	case schemaTypeInterface, schemaTypeObject:
		return "interface{}"
	case schemaTypeInt:
		return "int"
	case schemaTypeBoolean:
		return "bool"
	case schemaTypeString, schemaTypeMultiple:
		return "string"
	}

	return s
}

//func createChannels(m schemaPrefixList) *map[string]chan interface{} {
//	chans := make(map[string]chan interface{}, len(m))
//
//	for k := range m {
//		chans[k] = make(chan interface{}, 10)
//	}
//
//	return &chans
//}

func createByteChannels(m map[string]struct{}) map[string]chan []byte {
	chans := make(map[string]chan []byte, len(m))

	for k := range m {
		chans[k] = make(chan []byte, 10)
	}

	return chans
}

func checkMImports(items []IMethodItem, prefix string) bool {
	for _, v := range items {
		if IsNumber(v) {
			return true
		}

		if (IsBuiltin(v) || IsArray(v)) && strings.Count(v.GetGoType(), prefix) > 0 {
			return true
		}
	}

	return false
}

func checkTImports(item schemaJSONProperty, prefix string) bool {
	if IsBuiltin(item) && strings.Count(item.GetGoType(), prefix) > 0 {
		return true
	}

	if IsArray(item) && strings.Count(item.Items.GetGoType(), prefix) > 0 {
		return true
	}

	if IsObject(item) {
		for _, v := range item.GetProperties(false) {
			if (IsBuiltin(v) || IsArray(v) || IsNumber(v)) && strings.Count(v.GetGoType(), prefix) > 0 {
				return true
			} else if IsObject(v) {
				return checkTImports(v, prefix)
			}
		}
	}

	return false
}

func checkNames(tName, btName string) bool {
	btName = strings.Trim(btName, "[]()")

	return tName == btName
}

func checkChars(s string, chars string) bool {
	return strings.Count(s, chars) > 0
}
