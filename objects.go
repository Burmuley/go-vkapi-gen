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

type objectsSchema struct {
    keys        []string
    keyIndex    int
    initialized bool
    imports     map[string]map[string]struct{}
    Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func (o *objectsSchema) GetKey() string {
    return o.keys[o.keyIndex-1]
}

func (o *objectsSchema) Next() (IRender, bool) {
    if !o.initialized {
        o.keyIndex = 0
        o.initialized = true
    }

    if o.keyIndex < len(o.keys) {
        item := o.getItem()
        o.keyIndex++
        return item, true
    }

    return nil, false
}

func (o *objectsSchema) getItem() IRender {
    od := typeDefinition{}
    od[o.keys[o.keyIndex]] = o.Definitions[o.keys[o.keyIndex]]

    return &od
}

func (o *objectsSchema) Generate(outputDir string) error {
    tmplFuncs := make(map[string]interface{})
    tmplFuncs["checkNames"] = checkNames
    tmplFuncs = fillFuncs(tmplFuncs)

    tmplFuncs["deco"] = func(tName IType, rootType string) struct {
        T IType
        R string
    } {
        return struct {
            T IType
            R string
        }{T: tName, R: rootType}
    }

    _, tmplName := path.Split(objTmplName)

    tmpl, err := template.New(tmplName).Funcs(tmplFuncs).ParseFiles(objTmplName)

    if err != nil {
        return err
    }

    _, hTmplName := path.Split(objHeaderTmplName)

    hTmpl, err := template.New(hTmplName).Funcs(tmplFuncs).ParseFiles(objHeaderTmplName)

    if err != nil {
        return err
    }

    prefixes := map[string]struct{}{}

    for _, k := range o.keys {
        prefixes[getApiNamePrefix(k)] = struct{}{}
    }

    generateItems(o, hTmpl, tmpl, "objects", prefixes, o.imports)

    return nil
}

func (o *objectsSchema) Parse(fPath string) error {
    objects, err := loadSchemaFile(fPath)

    if err != nil {
        return fmt.Errorf("schema load error: %s", err)
    }

    logInfo(fmt.Sprintf("Successfully loaded schema from '%s'", fPath))

    if err := json.Unmarshal(objects, o); err != nil {
        return fmt.Errorf("JSON Error: %s", err)
    }

    objectsGlobal = o

    // fill the `stripPrefix` variable with 'true' for objects
    o.imports = make(map[string]map[string]struct{})

    for k := range o.Definitions {
        o.keys = append(o.keys, k)
        tmp := o.Definitions[k]
        setStripPrefix(&tmp, true)

        if tmp.GetType() == schemaTypeMultiple {
            if err := fillMultitype(&tmp, o); err != nil {
                return err
            }
        }

        o.Definitions[k] = tmp

        if checkTImports(tmp, "objects.") {
            //o.imports[getApiNamePrefix(k)][objectsImportPath] = struct{}{}
            addImport(o.imports, getApiNamePrefix(k), objectsImportPath)
        }

        if checkTImports(tmp, "responses.") {
            //o.imports[getApiNamePrefix(k)][responsesImportPath] = struct{}{}
            addImport(o.imports, getApiNamePrefix(k), responsesImportPath)
        }

        if checkTImports(tmp, "json.Number") {
            //o.imports[getApiNamePrefix(k)]["encoding/json"] = struct{}{}
            addImport(o.imports, getApiNamePrefix(k), "encoding/json")
        }
    }

    sort.Strings(o.keys)

    return nil
}
