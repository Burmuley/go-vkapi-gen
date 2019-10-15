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
	generateTypes(o.Definitions, OBJ_DIR_NAME, OBJ_HEADER_TMPL_NAME, OBJ_TMPL_NAME)

	return nil
}

func (o *objectsSchema) GetWriter() func() {
	return func() { return }
}

func (o *objectsSchema) Parse(fPath string) error {
	objects, err := loadSchemaFile(fPath)

	if err != nil {
		return fmt.Errorf("schema load error: %s", err)
	}

	//jsonObjects := objectsSchema{}

	if err := json.Unmarshal(objects, o); err != nil {
		return fmt.Errorf("JSON Error: %s", err)
	}

	return nil
}
