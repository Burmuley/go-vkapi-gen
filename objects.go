package main

import (
	"fmt"
	"log"
	"sync"
)

type objectsSchema struct {
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

func objectWriter(wg *sync.WaitGroup, ch chan map[string]schemaTyperChecker, filePrefix string) {
	schemaWriter(wg, ch, filePrefix, OBJ_DIR_NAME, OBJ_HEADER_TMPL_NAME, OBJ_TMPL_NAME)
}

func generateObjects(objects objectsSchema) {
	defCats := make(map[string]struct{})

	for k := range objects.Definitions {
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
		go objectWriter(wg, chans[k], k)
	}

	// Scan objects.Definitions and distribute data among appropriate channels
	for k, v := range objects.Definitions {
		tmp := make(map[string]schemaTyperChecker)
		tmp[k] = v

		if ch, ok := chans[getApiNamePrefix(k)]; ok {
			ch <- tmp
		} else {
			log.Fatal(fmt.Sprintf("channel '%s' not found in channels list", k))
		}
	}

	// Close all channels
	for _, v := range chans {
		close(v)
	}

	wg.Wait()
}
