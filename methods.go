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
	"log"
	"sync"
)

func methodsWriter(wg *sync.WaitGroup, ch chan IMethod, filePrefix string) {
	schemaMethodWriter(wg, ch, filePrefix, METHODS_HEADER_TMPL_NAME, METHODS_TMPL_NAME)
}

func createChannels(chList map[string]struct{}) (res *map[string]chan map[string]schemaTyperChecker) {
	// Create channels map and fill it
	chans := make(map[string]chan map[string]schemaTyperChecker, len(chList))

	for k := range chList {
		chans[k] = make(chan map[string]schemaTyperChecker, 10)
	}

	return &chans
}

func generateMethods(methods []IMethod) {
	logStep("Generating VK API methods")
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
