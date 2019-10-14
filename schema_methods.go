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
	"fmt"
)

type schemaApiError struct {
	Name  string `json:"name"`
	Code  int    `json:"code"`
	Descr string `json:"description"`
}

type schemaMethod struct {
	Name         string              `json:"name"`
	Descr        string              `json:"description"`
	AccessTokens []string            `json:"access_token_type"`
	Params       []*schemaMethodItem `json:"parameters"`
	Responses    struct {
		Response    *schemaMethodItem `json:"response"`
		ExtResponse *schemaMethodItem `json:"extendedResponse"`
	} `json:"responses"`
	Errors []*schemaApiError
}

func (s schemaMethod) GetDescription() string {
	if len(s.Descr) == 0 {
		return "NO DESCRIPTION IN JSON SCHEMA"
	}

	return s.Descr
}

func (s schemaMethod) GetResponses() []IMethodItem {
	tmp := make([]IMethodItem, 0)

	if s.Responses.Response != nil {
		tmp = append(tmp, s.Responses.Response)
	}

	if s.Responses.ExtResponse != nil {
		tmp = append(tmp, s.Responses.ExtResponse)
	}

	return tmp
}

func (s schemaMethod) GetResponse() IMethodItem {
	return s.Responses.Response
}

func (s schemaMethod) GetExtResponse() IMethodItem {
	return s.Responses.ExtResponse
}

func (s schemaMethod) GetParameters() []IMethodItem {
	mi := make([]IMethodItem, len(s.Params))

	for k, v := range s.Params {
		mi[k] = v
	}

	return mi
}

func (s schemaMethod) GetName() string {
	return s.Name
}

func (s schemaMethod) IsExtended() bool {
	return s.Responses.ExtResponse != nil
}

// Data structure implements method parameter
type schemaMethodItem struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Descr     string            `json:"description"`
	Required  bool              `json:"required"`
	Enum      []interface{}     `json:"enum"`
	EnumNames []string          `json:"enumNames"`
	Items     *schemaMethodItem `json:"items"`
	Ref       string            `json:"$ref"`
}

func (s schemaMethodItem) IsString() bool {
	return s.Type == SCHEMA_TYPE_STRING
}

func (s schemaMethodItem) IsInt() bool {
	return s.Type == SCHEMA_TYPE_INT
}

func (s schemaMethodItem) IsBuiltin() bool {
	return len(s.Ref) > 0
}

func (s schemaMethodItem) IsArray() bool {
	return s.Type == SCHEMA_TYPE_ARRAY
}

func (s schemaMethodItem) IsObject() bool {
	return s.Type == SCHEMA_TYPE_OBJECT
}

func (s schemaMethodItem) IsBoolean() bool {
	return s.Type == SCHEMA_TYPE_BOOLEAN
}

func (s schemaMethodItem) IsInterface() bool {
	return false
}

func (s schemaMethodItem) IsNumber() bool {
	return s.Type == SCHEMA_TYPE_NUMBER
}

func (s schemaMethodItem) IsMultiple() bool {
	return false
}

func (s schemaMethodItem) GetGoType() string {
	if len(s.Ref) > 0 {
		return getObjectTypeName(s.Ref)
	}

	if fmt.Sprint(s.Type) == SCHEMA_TYPE_ARRAY {
		return fmt.Sprintf("[]%s", s.Items.GetGoType())
	}

	return detectGoType(s.Type)
}

func (s schemaMethodItem) IsRequired() bool {
	return s.Required
}

func (s schemaMethodItem) GetType() string {
	if len(s.Ref) > 0 {
		return SCHEMA_TYPE_BUILTIN
	}

	return s.Type
}

func (s schemaMethodItem) GetName() string {
	return s.Name
}

func (s schemaMethodItem) GetDescription() string {
	if len(s.Descr) == 0 {
		return "NO DESCRIPTION IN JSON SCHEMA"
	}

	return s.Descr
}
