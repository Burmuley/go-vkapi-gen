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
	"encoding/json"
	"fmt"
	"strings"
)

type templateImports struct {
	Imports map[string]struct{}
	Prefix  string
}

type schemaPrefixList map[string]templateImports

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
		s.Type = schemaTypeMultiple
	default:
		s.Type = schemaTypeUnknown
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
	if len(s.AllOf) > 0 || len(s.OneOf) > 0 {
		return schemaTypeMultiple
	} else if len(s.Ref) > 0 {
		return schemaTypeBuiltin
	}

	if len(fmt.Sprint(s.Type)) > 0 {
		return fmt.Sprint(s.Type)
	} else if len(s.Properties) > 0 {
		return schemaTypeObject
	}

	return schemaTypeUnknown
}

func (s schemaJSONProperty) GetGoType(stripPrefix bool) (goTypes []string) {
	if s.AllOf != nil {
		for _, r := range s.AllOf {
			goTypes = append(goTypes, r.GetGoType(stripPrefix)...)
		}
		return
	} else if s.OneOf != nil {
		for _, r := range s.OneOf {
			goTypes = append(goTypes, r.GetGoType(stripPrefix)...)
		}
		return
	}

	if len(s.Ref) > 0 {
		var ref string

		if stripPrefix {
			stripped := strings.Split(s.Ref, "#")
			ref = strings.Join([]string{"#", stripped[len(stripped)-1]}, "")
		} else {
			ref = s.Ref
		}

		goTypes = append(goTypes, getObjectTypeName(ref))
		return
	}

	if fmt.Sprint(s.Type) == schemaTypeArray {
		return s.Items.GetGoType(stripPrefix)
	}

	goTypes = append(goTypes, detectGoType(fmt.Sprintf("%s", s.Type)))
	return
}

func (s schemaJSONProperty) GetDescription() string {
	return s.Descr
}

func (s schemaJSONProperty) GetProperties(stripPrefix bool) (pMap map[string]schemaJSONProperty) {
	if len(s.AllOf) > 0 || len(s.OneOf) > 0 {
		var mTypes []*schemaJSONProperty

		if len(s.AllOf) > 0 {
			mTypes = s.AllOf
		} else if len(s.OneOf) > 0 {
			mTypes = s.OneOf
		}

		pMap = make(map[string]schemaJSONProperty)

		for _, v := range mTypes {
			if v.IsBuiltin() {
				objName := convertName(strings.TrimLeft(v.GetGoType(stripPrefix)[0], "*"))

				pMap[objName] = *v
			} else if v.IsObject() {
				for k, v := range v.GetProperties(stripPrefix) {
					pMap[k] = v
				}
			}
		}

		return pMap
	}

	if len(s.Properties) > 0 {
		pMap = make(map[string]schemaJSONProperty, len(s.Properties))

		for k, v := range s.Properties {
			pMap[k] = *v
		}

		return pMap
	}

	return nil
}

func (s schemaJSONProperty) IsString() bool {
	return s.GetType() == schemaTypeString
}

func (s schemaJSONProperty) IsInt() bool {
	return s.GetType() == schemaTypeInt
}

func (s schemaJSONProperty) IsBuiltin() bool {
	return s.GetType() == schemaTypeBuiltin
}

func (s schemaJSONProperty) IsArray() bool {
	return s.GetType() == schemaTypeArray
}

func (s schemaJSONProperty) IsObject() bool {
	return s.GetType() == schemaTypeObject
}

func (s schemaJSONProperty) IsBoolean() bool {
	return s.GetType() == schemaTypeBoolean
}

func (s schemaJSONProperty) IsInterface() bool {
	return s.GetType() == schemaTypeInterface
}

func (s schemaJSONProperty) IsNumber() bool {
	return s.GetType() == schemaTypeNumber
}

func (s schemaJSONProperty) IsMultiple() bool {
	return s.GetType() == schemaTypeMultiple
}
