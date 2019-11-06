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
	"path"
	"sort"
	"text/template"
)

type responsesSchema struct {
	keys        []string
	keyIndex    int
	initialized bool
	imports     map[string]map[string]struct{}
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func (r *responsesSchema) Next() (IRender, bool) {
	if !r.initialized {
		r.keyIndex = 0
		r.initialized = true
	}

	if r.keyIndex < len(r.keys) {
		item := r.getItem()
		r.keyIndex++
		return item, true
	}

	return nil, false
}

func (r *responsesSchema) GetKey() string {
	return r.keys[r.keyIndex-1]
}

func (r *responsesSchema) getItem() IRender {
	od := typeDefinition{}
	od[r.keys[r.keyIndex]] = r.Definitions[r.keys[r.keyIndex]]

	return &od
}

func (r *responsesSchema) Generate(outputDir string) error {
	tmplFuncs := make(map[string]interface{})
	tmplFuncs["cutSuffix"] = cutSuffix
	tmplFuncs["checkNames"] = checkNames
	tmplFuncs = fillFuncs(tmplFuncs)

	//generateTypes(r.Definitions, outputDir, respDirName, respHeaderTmplName, respTmplName, tmplFuncs)

	_, tmplName := path.Split(respTmplName)

	tmpl, err := template.New(tmplName).Funcs(tmplFuncs).ParseFiles(respTmplName)

	if err != nil {
		return err
	}

	_, hTmplName := path.Split(respHeaderTmplName)

	hTmpl, err := template.New(hTmplName).Funcs(tmplFuncs).ParseFiles(respHeaderTmplName)

	if err != nil {
		return err
	}

	prefixes := map[string]struct{}{}

	for _, k := range r.keys {
		prefixes[getApiNamePrefix(k)] = struct{}{}
	}

	generateItems(r, hTmpl, tmpl, "responses", prefixes, r.imports)

	return nil
}

func (r *responsesSchema) Parse(fPath string) error {
	responses, err := loadSchemaFile(fPath)

	if err != nil {
		return fmt.Errorf("schema load error: %s", err)
	}

	if err := json.Unmarshal(responses, r); err != nil {
		return fmt.Errorf("JSON Error: %s", err)
	}

	for k := range r.Definitions {
		r.keys = append(r.keys, k)
	}

	sort.Strings(r.keys)

	return nil
}
