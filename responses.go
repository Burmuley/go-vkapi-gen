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

var schemaTypes = map[string]func() interface{}{
	"string":    func() interface{} { return &schemaPrimitive{} },
	"integer":   func() interface{} { return &schemaPrimitive{} },
	"boolean":   func() interface{} { return &schemaPrimitive{} },
	"interface": func() interface{} { return &schemaPrimitive{} },
	"array":     func() interface{} { return &schemaArray{} },
	"object":    func() interface{} { return &schemaObject{} },
	"builtin":   func() interface{} { return &schemaPrimitive{} },
}

func parseResponseJSON(b []byte, wrapper interface{}) error {
	tmp := &struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(b, tmp); err != nil {
		return err
	}

	if tmp.Type != "" {
		resp := schemaTypes[tmp.Type]()
		if err := json.Unmarshal(b, resp); err != nil {
			return err
		}

		*wrapper.(*interface{}) = resp
	} else {
		resp := schemaTypes["builtin"]()
		if err := json.Unmarshal(b, resp); err != nil {
			return err
		}

		*wrapper.(*interface{}) = resp
	}

	return nil
}

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

// schemaPrimitive used to represent the following entites:
// Response:
//  - integer
//  - string
// Array response items
// Response property fields
type schemaPrimitive struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

func (s schemaPrimitive) GetGoType() string {
	switch s.Type {
	case R_TYPE_BUILTIN:
		return s.Description
	case R_TYPE_INT:
		return "int"
	case R_TYPE_STRING:
		return "string"
	case R_TYPE_BOOLEAN:
		return "bool"
	case R_TYPE_INTERFACE:
		return "interface{}"
	}

	return "UNKNOWN"
}

func (s schemaPrimitive) GetDescription() string {
	if s.Type != R_TYPE_BUILTIN {
		return s.Description
	}

	return ""
}

func (s schemaPrimitive) IsString() bool {
	return s.Type == R_TYPE_STRING
}

func (s schemaPrimitive) IsInt() bool {
	return s.Type == R_TYPE_INT
}

func (s schemaPrimitive) IsBuiltin() bool {
	return s.Type == R_TYPE_BUILTIN
}

func (s schemaPrimitive) IsArray() bool {
	return false
}

func (s schemaPrimitive) IsObject() bool {
	return false
}

func (s schemaPrimitive) IsBoolean() bool {
	return s.Type == R_TYPE_BOOLEAN
}

func (s schemaPrimitive) IsInterface() bool {
	return s.Type == R_TYPE_INTERFACE
}

func (s *schemaPrimitive) UnmarshalJSON(b []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	} else {
		switch tmp := tmp.(type) {
		case map[string]interface{}:
			if v, ok := tmp["type"]; ok {
				s.Type = fmt.Sprintf("%s", v)
				s.Description = fmt.Sprintf("%s", tmp["description"])
			} else if v, ok := tmp["$ref"]; ok {
				s.Type = R_TYPE_BUILTIN
				s.Description = getObjectTypeName(fmt.Sprintf("%s", v))
			}
		case []interface{}:
			s.Type = R_TYPE_INTERFACE
		default:
			return fmt.Errorf("unable to determine type of '%s'", string(b))
		}
	}

	return nil
}

type schemaArray struct {
	Type  string          `json:"type"`
	Items schemaPrimitive `json:"items"`
}

func (s schemaArray) GetGoType() string {
	return strings.Join([]string{"[]", s.Items.GetGoType()}, "")
}

type schemaObject struct {
	Type       string                   `json:"type"`
	Properties map[string]schemaWrapper `json:"properties"`
}

type schemaWrapper struct {
	SchemaType      int
	SchemaArray     schemaArray
	SchemaPrimitive schemaPrimitive
}

func (s schemaWrapper) GetGoType() string {
	switch s.SchemaType {
	case TYPE_ARRAY:
		return s.SchemaArray.GetGoType()
	case TYPE_INT, TYPE_BUILTIN, TYPE_STRING, TYPE_BOOLEAN, TYPE_INTERFACE:
		return s.SchemaPrimitive.GetGoType()
	}

	return "UNKNOWN"
}

func (s schemaWrapper) GetDescription() string {
	switch s.SchemaType {
	case TYPE_INT, TYPE_STRING:
		return s.SchemaPrimitive.GetDescription()
	}

	return ""
}

func (s schemaWrapper) IsString() bool {
	return s.SchemaType == TYPE_STRING
}

func (s schemaWrapper) IsInt() bool {
	return s.SchemaType == TYPE_INT
}

func (s schemaWrapper) IsBuiltin() bool {
	return s.SchemaType == TYPE_BUILTIN
}

func (s schemaWrapper) IsArray() bool {
	return s.SchemaType == TYPE_ARRAY
}

func (s schemaWrapper) IsObject() bool {
	return false
}

func (s schemaWrapper) IsBoolean() bool {
	return s.SchemaType == TYPE_BOOLEAN
}

func (s schemaWrapper) IsInterface() bool {
	return s.SchemaType == TYPE_INTERFACE
}

func (s *schemaWrapper) UnmarshalJSON(rawData []byte) error {
	var tmp interface{}

	if err := parseResponseJSON(rawData, &tmp); err == nil {
		switch tmp := tmp.(type) {
		case *schemaPrimitive:
			switch tmp.Type {
			case R_TYPE_STRING:
				s.SchemaType = TYPE_STRING
				s.SchemaPrimitive = *tmp
			case R_TYPE_INT:
				s.SchemaType = TYPE_INT
				s.SchemaPrimitive = *tmp
			case R_TYPE_BOOLEAN:
				s.SchemaType = TYPE_BOOLEAN
				s.SchemaPrimitive = *tmp
			case R_TYPE_BUILTIN:
				s.SchemaType = TYPE_BUILTIN
				s.SchemaPrimitive = *tmp
			}
		case *schemaArray:
			s.SchemaType = TYPE_ARRAY
			s.SchemaArray = *tmp
		}
	} else {
		return err
	}

	return nil
}

// Common response wrapper to catch and process all possible entities
type propertyWrapper struct {
	PropertyType      int
	PropertyPrimitive schemaPrimitive
	PropertyObject    schemaObject
	PropertyArray     schemaArray
}

func (r propertyWrapper) GetGoType() string {
	switch r.PropertyType {
	case TYPE_ARRAY:
		return r.PropertyArray.GetGoType()
	case TYPE_INT, TYPE_BUILTIN, TYPE_STRING, TYPE_BOOLEAN:
		return r.PropertyPrimitive.GetGoType()
	}

	return "UNKNOWN"
}

func (r propertyWrapper) GetDescription() string {
	panic("implement me")
}

func (r propertyWrapper) IsString() bool {
	return r.PropertyType == TYPE_STRING
}

func (r propertyWrapper) IsInt() bool {
	return r.PropertyType == TYPE_INT
}

func (r propertyWrapper) IsBuiltin() bool {
	return r.PropertyType == TYPE_BUILTIN
}

func (r propertyWrapper) IsArray() bool {
	return r.PropertyType == TYPE_ARRAY
}

func (r propertyWrapper) IsObject() bool {
	return r.PropertyType == TYPE_OBJECT
}

func (r propertyWrapper) IsBoolean() bool {
	return r.PropertyType == TYPE_BOOLEAN
}

func (r propertyWrapper) IsInterface() bool {
	return r.PropertyType == TYPE_INTERFACE
}

func (r *propertyWrapper) UnmarshalJSON(rawData []byte) error {
	var tmp interface{}

	if err := parseResponseJSON(rawData, &tmp); err == nil {
		switch tmp := tmp.(type) {
		case *schemaPrimitive:
			switch tmp.Type {
			case R_TYPE_STRING:
				r.PropertyType = TYPE_STRING
			case R_TYPE_INT:
				r.PropertyType = TYPE_INT
			case R_TYPE_BUILTIN:
				r.PropertyType = TYPE_BUILTIN
			case R_TYPE_INTERFACE:
				r.PropertyType = TYPE_INTERFACE
			}

			r.PropertyPrimitive = *tmp
		case *schemaArray:
			r.PropertyType = TYPE_ARRAY
			r.PropertyArray = *tmp
		//case *schemaBuiltin:
		//	r.PropertyType = TYPE_BUILTIN
		//	r.PropertyBuiltin = *tmp
		case *schemaObject:
			r.PropertyType = TYPE_OBJECT
			r.PropertyObject = *tmp
		}
	} else {
		return err
	}

	return nil
}

func responseWriter(wg *sync.WaitGroup, ch chan map[string]responsesDefinition, filePrefix string) {
	var (
		f   *os.File
		err error
	)

	// Open new file
	if f, err = os.OpenFile(
		filepath.Join(OUTPUT_DIR_NAME, RESP_DIR_NAME, fmt.Sprintf("%s.go", filePrefix)),
		os.O_CREATE|os.O_RDWR,
		0644); err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()
	defer wg.Done()

	// Render header and write to the file
	tmpl, err := template.New(strings.Split(RESP_HEADER_TMPL_NAME, "/")[1]).ParseFiles(RESP_HEADER_TMPL_NAME)
	err = tmpl.Execute(f, filePrefix)

	// Read responses definitions from channel and append to the file
	funcs := make(map[string]interface{})
	funcs["convertName"] = convertName

	for {
		d, more := <-ch

		if more {
			tmpl, err := template.New(strings.Split(RESP_TMPL_NAME, "/")[1]).Funcs(funcs).ParseFiles(RESP_TMPL_NAME)

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

func generateResponses(responses responsesSchema) {
	//fmt.Printf("%#v\n", responses.Definitions)

	defCats := make(map[string]struct{})

	for k := range responses.Definitions {
		if _, ok := defCats[getApiNamePrefix(k)]; !ok {
			defCats[getApiNamePrefix(k)] = struct{}{}
		}
	}

	// Create channels map and fill it
	chans := make(map[string]chan map[string]responsesDefinition)
	wg := &sync.WaitGroup{}
	wg.Add(len(defCats))

	for k := range defCats {
		chans[k] = make(chan map[string]responsesDefinition, 10)
		go responseWriter(wg, chans[k], k)
	}

	// Scan responses.Definitions and distribute data among appropriate channels
	for k, v := range responses.Definitions {
		tmp := make(map[string]responsesDefinition)
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
