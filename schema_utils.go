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
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

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
				log.Printf("[[%s]] error formatting code: %s. Writing code as is...", fName, err)
				if n, e := f.Write(bb); e != nil {
					log.Printf("error writing %s: %s", fName, e)
				} else {
					log.Printf("successfully written %d bytes (unformatted) to %s.", n, fName)
				}
			} else {
				if n, e := f.Write(fmtCode); e != nil {
					log.Printf("error writing %s: %s", fName, e)
				} else {
					log.Printf("successfully written %d bytes to %s", n, fName)
				}

			}

			return
		}
	}
}

func schemaMethodWriter(wg *sync.WaitGroup, ch chan map[string]IMethod, prefix, headerTmpl, bodyTmpl string) {
	var (
		f   *os.File
		err error
		buf bytes.Buffer
	)

	defer wg.Done()

	fName := filepath.Join(OUTPUT_DIR_NAME, fmt.Sprintf("%s.go", prefix))

	// Check if a target file exists and remove it if so
	if checkFileExists(fName) {
		if err := os.Remove(fName); err != nil {
			log.Fatal(fmt.Sprintf("file '%s' exists and can't be removed! Error: %s", fName, err))
			return
		}

		log.Printf("removed file: %s\n", fName)
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
				log.Printf("[[%s]] error formatting code: %s. Writing code as is...", fName, err)
				if n, e := f.Write(bb); e != nil {
					log.Printf("error writing %s: %s", fName, e)
				} else {
					log.Printf("successfully written %d bytes (unformatted) to %s.", n, fName)
				}
			} else {
				if n, e := f.Write(fmtCode); e != nil {
					log.Printf("error writing %s: %s", fName, e)
				} else {
					log.Printf("successfully written %d bytes to %s", n, fName)
				}

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
