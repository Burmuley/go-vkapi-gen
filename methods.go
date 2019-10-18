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
	"log"
	"sort"
	"strings"
	"sync"
)

type schemaMethods struct {
	Errors  []schemaApiError `json:"errors"`
	Methods []schemaMethod   `json:"methods"`
}

func (s *schemaMethods) Parse(fPath string) error {
	methods, err := loadSchemaFile(fPath)

	if err != nil {
		return fmt.Errorf("schema load error: %s", err)
	}

	if err := json.Unmarshal(methods, s); err != nil {
		return fmt.Errorf("JSON Error: %s", err)
	}

	return nil
}

func (s *schemaMethods) Generate(outputDir string) error {

	iMethods := make([]IMethod, 0)

	for _, v := range s.Methods {
		iMethods = append(iMethods, v)
	}

	generateMethods(iMethods)

	return nil
}

func checkImports(items []IMethodItem, prefix string) bool {
	for _, v := range items {
		if (v.IsBuiltin() || v.IsArray()) && strings.Count(v.GetGoType(), prefix) > 0 {
			return true
		}
	}

	return false
}

func generateMethods(methods []IMethod) {
	//methodsCats := make(map[string]struct{})
	methodsCats := make(schemaPrefixList)

	for k := range methods {
		mPref := getApiMethodNamePrefix(methods[k].GetName())

		if _, ok := methodsCats[mPref]; !ok {
			methodsCats[mPref] = templateImports{
				Imports: make(map[string]struct{}),
				Prefix:  mPref,
			}
		}

		// Inspect parameters and fill imports
		if checkImports(methods[k].GetParameters(), "objects.") {
			methodsCats[mPref].Imports[objectsImportPath] = struct{}{}
		}

		// Inspect responses and fill imports
		if checkImports(methods[k].GetResponses(), "responses.") {
			methodsCats[mPref].Imports[responsesImportPath] = struct{}{}
		}

		if checkImports(methods[k].GetResponses(), "objects.") {
			methodsCats[mPref].Imports[objectsImportPath] = struct{}{}
		}
	}

	// Create channels map and fill it
	chans := *createChannels(methodsCats)

	wg := &sync.WaitGroup{}
	wg.Add(len(methodsCats))

	funcs := make(map[string]interface{})
	funcs["convertName"] = convertName
	funcs["convertParam"] = convertParam
	funcs["getMNameSuffix"] = getApiMethodNameSuffix
	funcs["getMNamePrefix"] = getApiMethodNamePrefix
	funcs["deco"] = func(method IMethod, count int) struct {
		M IMethod
		C int
	} {
		return struct {
			M IMethod
			C int
		}{M: method, C: count}
	}
	funcs["getFLetter"] = func(s string) string {
		return string(s[0])
	}

	for k := range methodsCats {
		go schemaWriter(wg, chans[k], methodsCats[k], k, "/", methodsHeaderTmplName, methodsTmplName, funcs)
	}

	//Scan methods and distribute data among appropriate channels
	sort.Slice(methods, func(i, j int) bool { return methods[i].GetName() < methods[j].GetName() })

	for _, v := range methods {
		if ch, ok := chans[getApiMethodNamePrefix(v.GetName())]; ok {
			ch <- v
		} else {
			log.Fatal(fmt.Sprintf("channel '%s' not found in channels list", getApiMethodNamePrefix(v.GetName())))
		}
	}

	// Close all channels
	for _, v := range chans {
		close(v)
	}

	wg.Wait()
}
