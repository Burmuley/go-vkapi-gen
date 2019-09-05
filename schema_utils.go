package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

func parseSchemaJSON(b []byte, wrapper interface{}) error {
	//tmp := &struct {
	//    Type string `json:"type"`
	//}{}
	tmp := &struct {
		Type interface{} `json:"type"`
	}{}

	if err := json.Unmarshal(b, tmp); err != nil {
		return err
	}

	if tmp.Type == nil {
		resp := schemaTypes["builtin"]()
		if err := json.Unmarshal(b, resp); err != nil {
			return err
		}

		*wrapper.(*interface{}) = resp
		return nil
	}

	switch tmp.Type.(type) {
	case []interface{}:
		resp := schemaTypes["string"]()
		if err := json.Unmarshal(b, resp); err != nil {
			return err
		}

		*wrapper.(*interface{}) = resp
	case string:
		if str := tmp.Type.(string); str != "" {
			resp := schemaTypes[str]()
			if err := json.Unmarshal(b, resp); err != nil {
				return err
			}

			*wrapper.(*interface{}) = resp
		}
	}

	return nil
}

func schemaWriter(wg *sync.WaitGroup, ch chan map[string]schemaTyperChecker, prefix, dir, headerTmpl, bodyTmpl string) {
	var (
		f   *os.File
		err error
	)

	// Open new file
	if f, err = os.OpenFile(
		filepath.Join(OUTPUT_DIR_NAME, dir, fmt.Sprintf("%s.go", prefix)),
		os.O_CREATE|os.O_RDWR|os.O_SYNC,
		0644); err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()
	defer wg.Done()

	// Render header and write to the file
	tmpl, err := template.New(strings.Split(headerTmpl, "/")[1]).ParseFiles(headerTmpl)
	err = tmpl.Execute(f, prefix)

	// Read responses definitions from channel and append to the file
	funcs := make(map[string]interface{})
	funcs["convertName"] = convertName

	for {
		d, more := <-ch

		if more {
			tmpl, err := template.New(strings.Split(bodyTmpl, "/")[1]).Funcs(funcs).ParseFiles(bodyTmpl)

			if err != nil {
				log.Fatal(err)
				return
			}

			err = tmpl.Execute(f, d)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			return
		}
	}
}
