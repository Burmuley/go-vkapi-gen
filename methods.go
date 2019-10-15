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
	"sync"
)

type schemaMethods struct {
	Errors  []schemaApiError `json:"errors"`
	Methods []schemaMethod   `json:"methods"`
}

func (s *schemaMethods) GetWriter() func() {
	panic("implement me")
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

func methodsWriter(wg *sync.WaitGroup, ch chan IMethod, filePrefix string) {
	schemaMethodWriter(wg, ch, filePrefix, methodsHeaderTmplName, methodsTmplName)
}

//func createChannels(chList map[string]struct{}) (res *map[string]chan map[string]ITypeChecker) {
//	// Create channels map and fill it
//	chans := make(map[string]chan map[string]ITypeChecker, len(chList))
//
//	for k := range chList {
//		chans[k] = make(chan map[string]ITypeChecker, 10)
//	}
//
//	return &chans
//}

func generateMethods(methods []IMethod) {
	methodsCats := make(map[string]struct{})

	for k := range methods {
		if _, ok := methodsCats[getApiMethodNamePrefix(methods[k].GetName())]; !ok {
			methodsCats[getApiMethodNamePrefix(methods[k].GetName())] = struct{}{}
		}
	}

	// Create channels map and fill it
	//chans := *createChannels(methodsCats)
	chans := make(map[string]chan IMethod, len(methodsCats))

	for k := range methodsCats {
		chans[k] = make(chan IMethod, 10)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(methodsCats))

	for k := range methodsCats {
		go methodsWriter(wg, chans[k], k)
	}

	//Scan methods and distribute data among appropriate channels
	sort.Slice(methods, func(i, j int) bool { return methods[i].GetName() < methods[j].GetName() })
	for _, v := range methods {
		var tmp IMethod
		tmp = v

		if ch, ok := chans[getApiMethodNamePrefix(v.GetName())]; ok {
			ch <- tmp
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
