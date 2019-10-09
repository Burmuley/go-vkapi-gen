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
	"sort"
	"sync"
)

type responsesSchema struct {
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func responseWriter(wg *sync.WaitGroup, ch chan map[string]schemaTyperChecker, filePrefix string) {
	schemaWriter(wg, ch, filePrefix, RESP_DIR_NAME, RESP_HEADER_TMPL_NAME, RESP_TMPL_NAME)
}

func generateResponses(responses responsesSchema) {
	logStep("Generating VK API responses")
	defCats := make(map[string]struct{})
	defKeys := make([]string, 0)

	for k := range responses.Definitions {
		defKeys = append(defKeys, k)
		if _, ok := defCats[getApiNamePrefix(k)]; !ok {
			defCats[getApiNamePrefix(k)] = struct{}{}
		}
	}

	// Create channels map and fill it
	chans := make(map[string]chan map[string]schemaTyperChecker)
	wg := &sync.WaitGroup{}
	wg.Add(len(defCats))

	for k := range defCats {
		chans[k] = make(chan map[string]schemaTyperChecker, 10)
		go responseWriter(wg, chans[k], k)
	}

	// Scan responses.Definitions and distribute data among appropriate channels
	sort.Strings(defKeys)
	for _, v := range defKeys {
		tmp := make(map[string]schemaTyperChecker)
		tmp[v] = responses.Definitions[v]

		if ch, ok := chans[getApiNamePrefix(v)]; ok {
			ch <- tmp
		} else {
			log.Fatal(fmt.Sprintf("channel '%s' not found in channels list", v))
		}
	}

	// Close all channels
	for _, v := range chans {
		close(v)
	}

	wg.Wait()
}
