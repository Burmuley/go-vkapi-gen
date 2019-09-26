package main

import (
	"sync"
)

func methodsWriter() {

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
	methodsCats := make(map[string]struct{})

	for k := range methods {
		if _, ok := methodsCats[getApiMethodNamePrefix(methods[k].GetName())]; !ok {
			methodsCats[getApiMethodNamePrefix(methods[k].GetName())] = struct{}{}
		}
	}

	// Create channels map and fill it
	chans := *createChannels(methodsCats)
	wg := &sync.WaitGroup{}
	wg.Add(len(methodsCats))

	for k := range methodsCats {
		go responseWriter(wg, chans[k], k)
	}

	// Scan responses.Definitions and distribute data among appropriate channels
	//for k, v := range methods {
	//	tmp := make(map[string]schemaTyperChecker)
	//	tmp[k] = v
	//
	//	if ch, ok := chans[getApiNamePrefix(k)]; ok {
	//		ch <- tmp
	//	} else {
	//		log.Fatal(fmt.Sprintf("channel '%s' not found in channels list", k))
	//	}
	//}

	// Close all channels
	for _, v := range chans {
		close(v)
	}

	wg.Wait()
}
