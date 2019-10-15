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

type responseDefinition map[string]schemaJSONProperty

func (r responseDefinition) GetPrefix() string {
	panic("implement me")
}

type responsesSchema struct {
	Definitions responseDefinition `json:"definitions"`
}

func (r *responsesSchema) Generate(outputDir string) error {
	tmplFuncs := make(map[string]interface{})
	tmplFuncs["convertName"] = convertName

	generateTypes(r.Definitions, outputDir, respDirName, respHeaderTmplName, respTmplName, tmplFuncs)

	return nil
}

func (r *responsesSchema) GetWriter() func() {
	return func() { return }
}

func (r *responsesSchema) Parse(fPath string) error {
	responses, err := loadSchemaFile(fPath)

	if err != nil {
		return fmt.Errorf("schema load error: %s", err)
	}

	if err := json.Unmarshal(responses, r); err != nil {
		return fmt.Errorf("JSON Error: %s", err)
	}

	return nil
}
