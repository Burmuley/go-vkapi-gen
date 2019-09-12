package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

func parseSchemaJSON(b []byte, wrapper interface{}) error {
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

func checkFileExists(f string) bool {
	finf, _ := os.Stat(f)
	return finf != nil
}

func schemaWriter(wg *sync.WaitGroup, ch chan map[string]schemaTyperChecker, prefix, dir, headerTmpl, bodyTmpl string) {
	var (
		f   *os.File
		err error
		buf bytes.Buffer
	)

	defer wg.Done()

	fName := filepath.Join(OUTPUT_DIR_NAME, dir, fmt.Sprintf("%s.go", prefix))

	// Check if a target file exists and remove it if so
	if checkFileExists(fName) {
		if err := os.Remove(fName); err != nil {
			log.Fatal(fmt.Sprintf("file '%s' exists and can't be removed! Error: %s", fName, err))
			return
		}

		log.Printf("removed file: %s", fName)
	}

	// Open new file
	if f, err = os.OpenFile(fName, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644); err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	// Render header and write to the file
	tmpl, err := template.New(strings.Split(headerTmpl, "/")[1]).ParseFiles(headerTmpl)
	err = tmpl.Execute(&buf, prefix)

	// Read responses definitions from channel and append to the file
	funcs := make(map[string]interface{})
	funcs["convertName"] = convertName

	for {
		d, more := <-ch

		if more {
			tmpl, err := template.New(strings.Split(bodyTmpl, "/")[1]).Funcs(funcs).ParseFiles(bodyTmpl)

			if err != nil {
				log.Println(err)
				return
			}

			err = tmpl.Execute(&buf, d)

			if err != nil {
				log.Println(err)
			}
		} else {
			bb := buf.Bytes()
			if fmtCode, err := format.Source(bb); err != nil {
				log.Printf("[[%s.go]] error formatting code: %s. Writing code as is...", prefix, err)
				f.Write(bb)
			} else {
				f.Write(fmtCode)
			}

			return
		}
	}
}

func detectGoType(s string) string {
	switch s {
	case SCHEMA_TYPE_NUMBER:
		return "float64"
	case SCHEMA_TYPE_INTERFACE:
		return "interface{}"
	case SCHEMA_TYPE_INT:
		return "int"
	case SCHEMA_TYPE_BOOLEAN:
		return "bool"
	case SCHEMA_TYPE_STRING:
		return "string"
	}

	return s
}
