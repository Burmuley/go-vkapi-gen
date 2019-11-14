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
	"sync"
	"text/template"
)

func bufWriter(wg *sync.WaitGroup, bCh, hCh chan []byte, prefix, outDir string) {
	var (
		f          *os.File
		err        error
		bBuf, hBuf bytes.Buffer
	)

	defer wg.Done()

	fName := filepath.Join(outputDirName, outDir, fmt.Sprintf("%s.go", prefix))

	// Check if a target file exists and remove it if so
	if checkFileExists(fName) {
		if err := os.Remove(fName); err != nil {
			logError(fmt.Errorf("file '%s' exists and can't be removed! Error: %s", fName, err))
			return
		}

		logInfo(fmt.Sprintf("removed file: %s", fName))
	}

	// Open new file
	if f, err = os.OpenFile(fName, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644); err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	// listen for body and header channels
	for {
		body, bOk := <-bCh
		bBuf.Write(body)

		if !bOk {
			break
		}
	}

	header := <-hCh

	// Add header on top of body buffer
	hBuf.Write(header)
	hBuf.Write(bBuf.Bytes())

	bb := hBuf.Bytes()

	// Format code && write to the file
	if fmtCode, err := format.Source(bb); err != nil {
		log.Printf("[[%s]] error formatting code: %s. Writing code as is...", fName, err)
		if n, e := f.Write(bb); e != nil {
			logError(fmt.Errorf("error writing %s: %s", fName, e))
		} else {
			logInfo(fmt.Sprintf("successfully written %d bytes (unformatted) to %s.", n, fName))
		}
	} else {
		if n, e := f.Write(fmtCode); e != nil {
			logError(fmt.Errorf("error writing %s: %s", fName, e))
		} else {
			logInfo(fmt.Sprintf("successfully written %d bytes to %s", n, fName))
		}

	}

}

func generateItems(items IIterator, hTmpl, bTmpl *template.Template, outDir string, prefixes map[string]struct{}, imports map[string]map[string]struct{}) {
	bodyChans := createByteChannels(prefixes)
	headChans := createByteChannels(prefixes)

	wg := &sync.WaitGroup{}
	wg.Add(len(bodyChans))

	for k := range prefixes {
		go bufWriter(wg, bodyChans[k], headChans[k], k, outDir)
	}

	for val, ok := items.Next(); ok; val, ok = items.Next() {
		if buf, err := val.Render(bTmpl); err == nil {
			bodyChans[getApiNamePrefix(items.GetKey())] <- buf
		} else {
			logError(err)
			logError(err.(template.ExecError).Err)
		}
	}

	for k := range prefixes {
		hBuf := bytes.Buffer{}
		close(bodyChans[k])
		tmp := templateImports{
			Imports: imports[k],
			Prefix:  k,
		}
		if err := hTmpl.Execute(&hBuf, tmp); err == nil {
			headChans[k] <- hBuf.Bytes()
		}
		close(headChans[k])
	}

	wg.Wait()
}

func IsString(t IType) bool {
	return t.GetType() == schemaTypeString
}

func IsInt(t IType) bool {
	return t.GetType() == schemaTypeInt
}

func IsBuiltin(t IType) bool {
	return t.GetType() == schemaTypeBuiltin
}

func IsArray(t IType) bool {
	return t.GetType() == schemaTypeArray
}

func IsObject(t IType) bool {
	return t.GetType() == schemaTypeObject
}

func IsBoolean(t IType) bool {
	return t.GetType() == schemaTypeBoolean
}

func IsInterface(t IType) bool {
	return t.GetType() == schemaTypeInterface
}

func IsNumber(t IType) bool {
	return t.GetType() == schemaTypeNumber
}

func IsMultiple(t IType) bool {
	return t.GetType() == schemaTypeMultiple
}

func fillFuncs(m map[string]interface{}) map[string]interface{} {
	m["IsString"] = IsString
	m["IsInt"] = IsInt
	m["IsBuiltin"] = IsBuiltin
	m["IsArray"] = IsArray
	m["IsObject"] = IsObject
	m["IsBoolean"] = IsBoolean
	m["IsInterface"] = IsInterface
	m["IsNumber"] = IsNumber
	m["IsMultiple"] = IsMultiple
	m["checkChars"] = checkChars
	m["convertName"] = convertName
	return m
}
