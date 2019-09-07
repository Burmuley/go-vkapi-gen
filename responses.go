package main

import (
	"fmt"
	"log"
	"sync"
)

type responsesSchema struct {
	Definitions map[string]responsesDefinition `json:"definitions"`
}

type responsesDefinition struct {
	Type       string `json:"type"`
	Properties struct {
		Response propertyWrapper `json:"response"`
	} `json:"properties"`
	Description string `json:"description,omitempty"`
}

func (r responsesDefinition) IsString() bool {
	return r.Type == SCHEMA_TYPE_STRING
}

func (r responsesDefinition) IsInt() bool {
	return r.Type == SCHEMA_TYPE_INT
}

func (r responsesDefinition) IsBuiltin() bool {
	return r.Type == SCHEMA_TYPE_BUILTIN
}

func (r responsesDefinition) IsArray() bool {
	return r.Type == SCHEMA_TYPE_ARRAY
}

func (r responsesDefinition) IsObject() bool {
	return r.Type == SCHEMA_TYPE_OBJECT
}

func (r responsesDefinition) IsBoolean() bool {
	return r.Type == SCHEMA_TYPE_BOOLEAN
}

func (r responsesDefinition) IsInterface() bool {
	return r.Type == SCHEMA_TYPE_INTERFACE
}

func (r responsesDefinition) GetGoType() string {
	return r.Type
}

func (r responsesDefinition) GetDescription() string {
	return r.Description
}

func responseWriter(wg *sync.WaitGroup, ch chan map[string]schemaTyperChecker, filePrefix string) {
	schemaWriter(wg, ch, filePrefix, RESP_DIR_NAME, RESP_HEADER_TMPL_NAME, RESP_TMPL_NAME)
}

func generateResponses(responses responsesSchema) {
	defCats := make(map[string]struct{})

	for k := range responses.Definitions {
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
	for k, v := range responses.Definitions {
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
