package main

import (
	"encoding/json"
	"fmt"
)

type schema struct {
	Title       string                        `json:"title"`
	Definitions map[string]schemaJSONProperty `json:"definitions"`
}

type schemaTypeWrapper struct {
	Type string `json:"-"`
}

func (s schemaTypeWrapper) String() string {
	return fmt.Sprintf("%s", s.Type)
}

func (s *schemaTypeWrapper) UnmarshalJSON(b []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp.(type) {
	case string:
		s.Type = fmt.Sprintf("%s", tmp)
	case []interface{}:
		s.Type = SCHEMA_TYPE_INTERFACE
	}

	return nil
}

type schemaItemsWrapper struct {
	Items    *schemaJSONProperty   `json:"-"`
	ItemsArr []*schemaJSONProperty `json:"-"`
}

func (s schemaItemsWrapper) GetGoType() (tmp []string) {
	if s.ItemsArr != nil {
		tmp = append(tmp, "interface{}")
		return tmp
	} else if s.Items != nil {
		for _, i := range s.Items.GetGoType() {
			tmp = append(tmp, fmt.Sprintf("%s", i))
		}
		return tmp
	}

	return []string{}
}

func (s schemaItemsWrapper) GetDescription() string {
	return ""
}

func (s *schemaItemsWrapper) UnmarshalJSON(b []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp.(type) {
	case []interface{}:
		err := json.Unmarshal(b, &s.ItemsArr)
		return err
	default:
		err := json.Unmarshal(b, &s.Items)
		return err
	}
}

type schemaJSONProperty struct {
	Type       schemaTypeWrapper              `json:"type,omitempty"`
	Descr      string                         `json:"description,omitempty"`
	AllOf      []*schemaJSONProperty          `json:"allOf,omitempty"`
	OneOf      []*schemaJSONProperty          `json:"oneOf,omitempty"`
	Properties map[string]*schemaJSONProperty `json:"properties,omitempty"`
	Required   []string                       `json:"required,omitempty"`
	Enum       []interface{}                  `json:"enum,omitempty"` // TODO: make a wrapper (can be int or string)
	EnumNames  []string                       `json:"enum_names,omitempty"`
	Items      *schemaItemsWrapper            `json:"items,omitempty"`
	Ref        string                         `json:"$ref,omitempty"`
}

func (s schemaJSONProperty) GetType() string {
	if len(fmt.Sprint(s.Type)) > 0 {
		return fmt.Sprint(s.Type)
	} else {
		if len(s.AllOf) > 0 || len(s.OneOf) > 0 || len(s.Properties) > 0 {
			return SCHEMA_TYPE_OBJECT
		} else if len(s.Ref) > 0 {
			return SCHEMA_TYPE_BUILTIN
		}
	}

	return "<<<UNKNOWN>>>"
}

func (s schemaJSONProperty) GetGoType() (tmp []string) {
	if s.AllOf != nil {
		for _, r := range s.AllOf {
			tmp = append(tmp, r.GetGoType()...)
		}
		return tmp
	}

	if len(s.Ref) > 0 {
		tmp = append(tmp, getObjectTypeName(s.Ref))
		return tmp
	}

	if fmt.Sprint(s.Type) == SCHEMA_TYPE_ARRAY {
		return s.Items.GetGoType()
	}

	tmp = append(tmp, detectGoType(fmt.Sprintf("%s", s.Type)))
	return tmp
}

func (s schemaJSONProperty) GetDescription() string {
	return s.Descr
}

func (s schemaJSONProperty) GetProperties() map[string]schemaJSONProperty {
	tmp := make(map[string]schemaJSONProperty, len(s.Properties))

	for k, v := range s.Properties {
		tmp[k] = *v
	}

	return tmp
}

func (s schemaJSONProperty) IsString() bool {
	return s.GetType() == SCHEMA_TYPE_STRING
}

func (s schemaJSONProperty) IsInt() bool {
	return s.GetType() == SCHEMA_TYPE_INT
}

func (s schemaJSONProperty) IsBuiltin() bool {
	return s.GetType() == SCHEMA_TYPE_BUILTIN
}

func (s schemaJSONProperty) IsArray() bool {
	return s.GetType() == SCHEMA_TYPE_ARRAY
}

func (s schemaJSONProperty) IsObject() bool {
	return s.GetType() == SCHEMA_TYPE_OBJECT
}

func (s schemaJSONProperty) IsBoolean() bool {
	return s.GetType() == SCHEMA_TYPE_BOOLEAN
}

func (s schemaJSONProperty) IsInterface() bool {
	return s.GetType() == SCHEMA_TYPE_INTERFACE
}

func (s schemaJSONProperty) IsNumber() bool {
	return s.GetType() == SCHEMA_TYPE_NUMBER
}
