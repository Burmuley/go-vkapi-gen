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

type schemaMethods struct {
	Errors  []schemaMethodsErrors  `json:"errors"`
	Methods []schemaMethodsMethods `json:"methods"`
}

type schemaMethodsErrors struct {
	Name  string `json:"name"`
	Code  int    `json:"code"`
	Descr string `json:"description"`
}

type schemaMethodsMethods struct {
	Name         string               `json:"name"`
	Descr        string               `json:"description"`
	AccessTokens []string             `json:"access_token_type"`
	Params       []*schemaMethodsItem `json:"parameters"`
	Responses    struct {
		Response    *schemaMethodsItem `json:"response"`
		ExtResponse *schemaMethodsItem `json:"extendedResponse"`
	} `json:"responses"`
	Errors []*schemaMethodsErrors
}

func (s schemaMethodsMethods) GetDescription() string {
	if len(s.Descr) == 0 {
		return "NO DESCRIPTION IN JSON SCHEMA"
	}

	return s.Descr
}

func (s schemaMethodsMethods) GetResponses() []IMethodItem {
	tmp := make([]IMethodItem, 0)

	if s.Responses.Response != nil {
		tmp = append(tmp, s.Responses.Response)
	}

	if s.Responses.ExtResponse != nil {
		tmp = append(tmp, s.Responses.ExtResponse)
	}

	return tmp
}

func (s schemaMethodsMethods) GetResponse() IMethodItem {
	return s.Responses.Response
}

func (s schemaMethodsMethods) GetExtResponse() IMethodItem {
	return s.Responses.ExtResponse
}

func (s schemaMethodsMethods) GetParameters() []IMethodItem {
	mi := make([]IMethodItem, len(s.Params))

	for k, v := range s.Params {
		mi[k] = v
	}

	return mi
}

func (s schemaMethodsMethods) GetName() string {
	return s.Name
}

func (s schemaMethodsMethods) IsExtended() bool {
	return s.Responses.ExtResponse != nil
}

// Data structure implements method parameter
type schemaMethodsItem struct {
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	Descr     string             `json:"description"`
	Required  bool               `json:"required"`
	Enum      []interface{}      `json:"enum"`
	EnumNames []string           `json:"enumNames"`
	Items     *schemaMethodsItem `json:"items"`
	Ref       string             `json:"$ref"`
}

func (s schemaMethodsItem) IsString() bool {
	return s.Type == SCHEMA_TYPE_STRING
}

func (s schemaMethodsItem) IsInt() bool {
	return s.Type == SCHEMA_TYPE_INT
}

func (s schemaMethodsItem) IsBuiltin() bool {
	return len(s.Ref) > 0
}

func (s schemaMethodsItem) IsArray() bool {
	return s.Type == SCHEMA_TYPE_ARRAY
}

func (s schemaMethodsItem) IsObject() bool {
	return s.Type == SCHEMA_TYPE_OBJECT
}

func (s schemaMethodsItem) IsBoolean() bool {
	return s.Type == SCHEMA_TYPE_BOOLEAN
}

func (s schemaMethodsItem) IsInterface() bool {
	return false
}

func (s schemaMethodsItem) IsNumber() bool {
	return s.Type == SCHEMA_TYPE_NUMBER
}

func (s schemaMethodsItem) IsMultiple() bool {
	return false
}

func (s schemaMethodsItem) GetGoType() string {
	if len(s.Ref) > 0 {
		return getObjectTypeName(s.Ref)
	}

	if fmt.Sprint(s.Type) == SCHEMA_TYPE_ARRAY {
		return s.Items.GetGoType()
	}

	return detectGoType(s.Type)

	//tmp = append(tmp, detectGoType(fmt.Sprintf("%s", s.Type)))
	//return tmp
}

func (s schemaMethodsItem) IsRequired() bool {
	return s.Required
}

func (s schemaMethodsItem) GetType() string {
	if len(s.Ref) > 0 {
		return SCHEMA_TYPE_BUILTIN
	}

	return s.Type
}

func (s schemaMethodsItem) GetName() string {
	return s.Name
}

func (s schemaMethodsItem) GetDescription() string {
	if len(s.Descr) == 0 {
		return "NO DESCRIPTION IN JSON SCHEMA"
	}

	return s.Descr
}
