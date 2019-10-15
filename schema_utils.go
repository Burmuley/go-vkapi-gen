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
	"sort"
	"strings"
	"sync"
	"text/template"
)

func schemaWriter(wg *sync.WaitGroup, ch chan interface{}, prefix, dir, headerTmpl, bodyTmpl string, tmplFuncs map[string]interface{}) {
	var (
		f   *os.File
		err error
		buf bytes.Buffer
	)

	defer wg.Done()

	fName := filepath.Join(outputDirName, dir, fmt.Sprintf("%s.go", prefix))

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
	tmpl, err := template.New(strings.Split(headerTmpl, "/")[1]).Funcs(tmplFuncs).ParseFiles(headerTmpl)
	err = tmpl.Execute(&buf, prefix)

	// Read responses definitions from channel and append to the file
	//tmplFuncs := make(map[string]interface{})
	//tmplFuncs["convertName"] = convertName

	for {
		d, more := <-ch

		if more {
			tmpl, err := template.New(strings.Split(bodyTmpl, "/")[1]).Funcs(tmplFuncs).ParseFiles(bodyTmpl)

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
	case schemaTypeNumber:
		return "float64"
	case schemaTypeInterface:
		return "interface{}"
	case schemaTypeInt:
		return "int"
	case schemaTypeBoolean:
		return "bool"
	case schemaTypeString:
		return "string"
	}

	return s
}

func generateTypes(types map[string]schemaJSONProperty, outRootDir, dir, headerTmpl, bodyTmpl string, tmplFuncs map[string]interface{}) {
	defCats := make(map[string]struct{})
	defKeys := make([]string, 0)

	for k := range types {
		defKeys = append(defKeys, k)
		if _, ok := defCats[getApiNamePrefix(k)]; !ok {
			defCats[getApiNamePrefix(k)] = struct{}{}
		}
	}

	// Create channels map and fill it
	chans := make(map[string]chan interface{})
	wg := &sync.WaitGroup{}
	wg.Add(len(defCats))

	for k := range defCats {
		chans[k] = make(chan interface{}, 10)
		go schemaWriter(wg, chans[k], k, dir, headerTmpl, bodyTmpl, tmplFuncs)
	}

	// Scan types and distribute data among appropriate channels
	sort.Strings(defKeys)
	for _, v := range defKeys {
		tmp := make(map[string]ITypeChecker)
		tmp[v] = types[v]

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
