package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var schemaTypes = map[string]func() interface{}{
	"string":    func() interface{} { return &schemaPrimitive{} },
	"integer":   func() interface{} { return &schemaPrimitive{} },
	"number":    func() interface{} { return &schemaPrimitive{} },
	"boolean":   func() interface{} { return &schemaPrimitive{} },
	"interface": func() interface{} { return &schemaPrimitive{} },
	"array":     func() interface{} { return &schemaArray{} },
	"object":    func() interface{} { return &schemaObject{} },
	"builtin":   func() interface{} { return &schemaPrimitive{} },
}

// schemaPrimitive used to represent the following entites:
// Response:
//  - integer
//  - string
// Array response items
// Response property fields
type schemaPrimitive struct {
	Type            string `json:"type"`
	Description     string `json:"description,omitempty"`
	BuiltinTypeName string `json:"-"`
}

func (s schemaPrimitive) GetGoType() string {
	switch s.Type {
	case SCHEMA_TYPE_BUILTIN:
		return s.BuiltinTypeName
	case SCHEMA_TYPE_INT:
		return "int"
	case SCHEMA_TYPE_STRING:
		return "string"
	case SCHEMA_TYPE_BOOLEAN:
		return "bool"
	case SCHEMA_TYPE_INTERFACE:
		return "interface{}"
	}

	return "UNKNOWN"
}

func (s schemaPrimitive) GetDescription() string {
	if len(s.Description) == 0 {
		return ""
	}

	return s.Description
}

func (s schemaPrimitive) IsString() bool {
	return s.Type == SCHEMA_TYPE_STRING
}

func (s schemaPrimitive) IsInt() bool {
	return s.Type == SCHEMA_TYPE_INT
}

func (s schemaPrimitive) IsBuiltin() bool {
	return s.Type == SCHEMA_TYPE_BUILTIN
}

func (s schemaPrimitive) IsArray() bool {
	return false
}

func (s schemaPrimitive) IsObject() bool {
	return false
}

func (s schemaPrimitive) IsBoolean() bool {
	return s.Type == SCHEMA_TYPE_BOOLEAN
}

func (s schemaPrimitive) IsInterface() bool {
	return s.Type == SCHEMA_TYPE_INTERFACE
}

func (s *schemaPrimitive) UnmarshalJSON(b []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp := tmp.(type) {
	case map[string]interface{}:
		if v, ok := tmp["type"]; ok {
			switch v.(type) {
			case []interface{}:
				s.Type = "string"
			case string:
				s.Type = fmt.Sprintf("%s", v)
			}
			s.Description = fmt.Sprintf("%s", tmp["description"])
		} else if v, ok := tmp["$ref"]; ok {
			s.Type = SCHEMA_TYPE_BUILTIN
			if dsc, ok := tmp["description"]; ok {
				s.Description = fmt.Sprintf("%s", dsc)
			}
			s.BuiltinTypeName = getObjectTypeName(fmt.Sprintf("%s", v))
		}
	case []interface{}:
		s.Type = SCHEMA_TYPE_INTERFACE
	default:
		return fmt.Errorf("unable to determine type of '%s'", string(b))
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
	case TYPE_INT, TYPE_STRING, TYPE_BOOLEAN, TYPE_BUILTIN:
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

	if err := parseSchemaJSON(rawData, &tmp); err == nil {
		switch tmp := tmp.(type) {
		case *schemaPrimitive:
			switch tmp.Type {
			case SCHEMA_TYPE_STRING:
				s.SchemaType = TYPE_STRING
				s.SchemaPrimitive = *tmp
			case SCHEMA_TYPE_INT:
				s.SchemaType = TYPE_INT
				s.SchemaPrimitive = *tmp
			case SCHEMA_TYPE_BOOLEAN:
				s.SchemaType = TYPE_BOOLEAN
				s.SchemaPrimitive = *tmp
			case SCHEMA_TYPE_BUILTIN:
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
	switch r.PropertyType {
	case TYPE_ARRAY:
		return ""
	case TYPE_INT, TYPE_BUILTIN, TYPE_STRING, TYPE_BOOLEAN:
		return r.PropertyPrimitive.GetDescription()
	}

	return "UNKNOWN"
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

	if err := parseSchemaJSON(rawData, &tmp); err == nil {
		switch tmp := tmp.(type) {
		case *schemaPrimitive:
			switch tmp.Type {
			case SCHEMA_TYPE_STRING:
				r.PropertyType = TYPE_STRING
			case SCHEMA_TYPE_INT:
				r.PropertyType = TYPE_INT
			case SCHEMA_TYPE_BUILTIN:
				r.PropertyType = TYPE_BUILTIN
			case SCHEMA_TYPE_INTERFACE:
				r.PropertyType = TYPE_INTERFACE
			}

			r.PropertyPrimitive = *tmp
		case *schemaArray:
			r.PropertyType = TYPE_ARRAY
			r.PropertyArray = *tmp
		case *schemaObject:
			r.PropertyType = TYPE_OBJECT
			r.PropertyObject = *tmp
		}
	} else {
		fmt.Printf("JSON ERR BODY: %s", string(rawData))
		return err
	}

	return nil
}
