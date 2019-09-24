package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

//////////////////////////////////////////////////////////////////////
// JSON schema main structure
//////////////////////////////////////////////////////////////////////
type schema struct {
	Title       string                         `json:"title"`
	Definitions map[string]*schemaJSONProperty `json:"definitions"`
}

//////////////////////////////////////////////////////////////////////
// JSON schema `type` field wrapper
//////////////////////////////////////////////////////////////////////
type schemaTypeWrapper struct {
	Type string `json:"-"`
}

func (s schemaTypeWrapper) String() string {
	return fmt.Sprintf("%s", s.Type)
}

func (s *schemaTypeWrapper) UnmarshalJSON(b []byte) error {
	var tmp interface{}

	if err := json.Unmarshal(b, &tmp); err != nil {
		return schemaError{string(b), err}
	}

	switch tmp.(type) {
	case string:
		s.Type = fmt.Sprintf("%s", tmp)
	case []interface{}:
		s.Type = SCHEMA_TYPE_MULTIPLE
	default:
		s.Type = SCHEMA_TYPE_UNKNOWN
		return schemaError{
			errInfo: string(b),
			err:     fmt.Errorf("%s", "unknown schema type"),
		}
	}

	return nil

}

//////////////////////////////////////////////////////////////////////
// JSON schema `items` field wrapper
//////////////////////////////////////////////////////////////////////
type schemaItemsWrapper struct {
	Items    *schemaJSONProperty   `json:"-"`
	ItemsArr []*schemaJSONProperty `json:"-"`
}

func (s schemaItemsWrapper) GetGoType(stripPrefix bool) (tmp []string) {
	if s.ItemsArr != nil {
		tmp = append(tmp, "interface{}")
		return tmp
	} else if s.Items != nil {
		for _, i := range s.Items.GetGoType(stripPrefix) {
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
		return schemaError{string(b), err}
	}

	switch tmp.(type) {
	case []interface{}:
		if err := json.Unmarshal(b, &s.ItemsArr); err != nil {
			return schemaError{string(b), err}
		}

		return nil
	default:
		if err := json.Unmarshal(b, &s.Items); err != nil {
			return schemaError{string(b), err}
		}

		return nil
	}
}

//////////////////////////////////////////////////////////////////////
// JSON schema element data structure
//////////////////////////////////////////////////////////////////////
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

	return SCHEMA_TYPE_UNKNOWN
}

func (s schemaJSONProperty) GetGoType(stripPrefix bool) (tmp []string) {
	if s.AllOf != nil {
		for _, r := range s.AllOf {
			tmp = append(tmp, r.GetGoType(stripPrefix)...)
		}
		return tmp
	}

	if len(s.Ref) > 0 {
		var ref string

		if stripPrefix {
			stripped := strings.Split(s.Ref, "#")
			ref = strings.Join([]string{"#", stripped[len(stripped)-1]}, "")
		} else {
			ref = s.Ref
		}

		tmp = append(tmp, getObjectTypeName(ref))
		return tmp
	}

	if fmt.Sprint(s.Type) == SCHEMA_TYPE_ARRAY {
		return s.Items.GetGoType(stripPrefix)
	}

	tmp = append(tmp, detectGoType(fmt.Sprintf("%s", s.Type)))
	return tmp
}

func (s schemaJSONProperty) GetDescription() string {
	return s.Descr
}

func (s schemaJSONProperty) GetProperties(stripPrefix bool) (tmp map[string]schemaJSONProperty) {
	if len(s.AllOf) > 0 || len(s.OneOf) > 0 {
		var tmpOf []*schemaJSONProperty

		if len(s.AllOf) > 0 {
			tmpOf = s.AllOf
		} else if len(s.OneOf) > 0 {
			tmpOf = s.OneOf
		}

		tmp = make(map[string]schemaJSONProperty, len(tmpOf))

		for _, v := range tmpOf {
			if v.IsBuiltin() {
				objName := convertName(v.GetGoType(stripPrefix)[0])

				tmp[objName] = *v
			} else if v.IsObject() {
				for k, v := range v.GetProperties(stripPrefix) {
					tmp[k] = v
				}
			}
		}

		return tmp
	}

	if len(s.Properties) > 0 {
		tmp = make(map[string]schemaJSONProperty, len(s.Properties))

		for k, v := range s.Properties {
			tmp[k] = *v
		}

		return tmp
	}

	return nil
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

func (s schemaJSONProperty) IsMultiple() bool {
	return s.GetType() == SCHEMA_TYPE_MULTIPLE
}
