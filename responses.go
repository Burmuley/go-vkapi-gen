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

type responsesSchema struct {
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func (r *responsesSchema) Generate(outputDir string) error {
	generateTypes(r.Definitions, RESP_DIR_NAME, RESP_HEADER_TMPL_NAME, RESP_TMPL_NAME)

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
