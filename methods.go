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
    "text/template"
)

type schemaMethods struct {
    keys        []string
    keyIndex    int
    initialized bool
    imports     map[string]map[string]struct{}
    Errors      []schemaApiError `json:"errors"`
    Methods     []schemaMethod   `json:"methods"`
}

func (s *schemaMethods) Next() (IRender, bool) {
    if !s.initialized {
        s.keyIndex = 0
        s.initialized = true
    }

    if s.keyIndex < len(s.keys) {
        item := s.getItem()
        s.keyIndex++
        return item, true
    }

    return nil, false
}

func (s *schemaMethods) GetKey() string {
    return s.keys[s.keyIndex-1]
}

func (s *schemaMethods) getItem() IRender {
    od := schemaMethod{}
    od = s.Methods[s.keyIndex]

    return &od
}

func (s *schemaMethods) Parse(fPath string) error {
    methods, err := loadSchemaFile(fPath)

    if err != nil {
        return fmt.Errorf("schema load error: %s", err)
    }

    if err := json.Unmarshal(methods, s); err != nil {
        return fmt.Errorf("JSON Error: %s", err)
    }

    s.imports = make(map[string]map[string]struct{})

    for k := range s.Methods {
        s.keys = append(s.keys, s.Methods[k].GetName())
        mPref := getApiNamePrefix(s.Methods[k].GetName())

        // Inspect parameters and fill imports
        if checkMImports(s.Methods[k].GetParameters(), "objects.") {
            //s.imports[mPref][objectsImportPath] = struct{}{}
            addImport(s.imports, mPref, objectsImportPath)
        }

        if checkMImports(s.Methods[k].GetParameters(), "json.Number") {
            //s.imports[mPref]["encoding/json"] = struct{}{}
            addImport(s.imports, mPref, "encoding/json")
        }

        // Inspect responses and fill imports
        if checkMImports(s.Methods[k].GetResponses(), "responses.") {
            //s.imports[mPref][responsesImportPath] = struct{}{}
            addImport(s.imports, mPref, responsesImportPath)
        }

        if checkMImports(s.Methods[k].GetResponses(), "objects.") {
            //s.imports[mPref][objectsImportPath] = struct{}{}
            addImport(s.imports, mPref, objectsImportPath)
        }
    }

    return nil
}

func (s *schemaMethods) Generate(outputDir string) error {
    tmplFuncs := make(map[string]interface{})
    tmplFuncs = fillFuncs(tmplFuncs)
    tmplFuncs["convertParam"] = convertParam
    tmplFuncs["getMNameSuffix"] = getApiMethodNameSuffix
    tmplFuncs["getMNamePrefix"] = getApiNamePrefix
    tmplFuncs["cutSuffix"] = cutSuffix
    tmplFuncs["deco"] = func(method IMethod, count int) struct {
        M IMethod
        C int
    } {
        return struct {
            M IMethod
            C int
        }{M: method, C: count}
    }
    tmplFuncs["getFLetter"] = func(s string) string {
        return string(s[0])
    }

    _, tmplName := path.Split(methodsTmplName)

    tmpl, err := template.New(tmplName).Funcs(tmplFuncs).ParseFiles(methodsTmplName)

    if err != nil {
        return err
    }

    _, hTmplName := path.Split(methodsHeaderTmplName)

    hTmpl, err := template.New(hTmplName).Funcs(tmplFuncs).ParseFiles(methodsHeaderTmplName)

    if err != nil {
        return err
    }

    prefixes := map[string]struct{}{}

    for _, v := range s.Methods {
        prefixes[getApiNamePrefix(v.GetName())] = struct{}{}
    }

    generateItems(s, hTmpl, tmpl, "/", prefixes, s.imports)

    return nil
}
