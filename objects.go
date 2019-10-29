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
)

type objectsSchema struct {
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func (o *objectsSchema) Generate(outputDir string) error {
	tmplFuncs := make(map[string]interface{})
	tmplFuncs["convertName"] = convertName
	tmplFuncs["checkNames"] = checkNames
	tmplFuncs = fillFuncs(tmplFuncs)

	tmplFuncs["deco"] = func(tName schemaJSONProperty, rootType string) struct {
		T schemaJSONProperty
		R string
	} {
		return struct {
			T schemaJSONProperty
			R string
		}{T: tName, R: rootType}
	}

	generateTypes(o.Definitions, outputDir, objDirName, objHeaderTmplName, objTmplName, tmplFuncs)

	return nil
}

func (o *objectsSchema) Parse(fPath string) error {
	objects, err := loadSchemaFile(fPath)

	if err != nil {
		return fmt.Errorf("schema load error: %s", err)
	}

	if err := json.Unmarshal(objects, o); err != nil {
		return fmt.Errorf("JSON Error: %s", err)
	}

	return nil
}
