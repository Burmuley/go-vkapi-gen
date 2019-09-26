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
	AccessTokens []string             `json:"access_token_type"`
	Params       []*schemaMethodsItem `json:"parameters"`
	Responses    struct {
		Response    *schemaMethodsItem `json:"response"`
		ExtResponse *schemaMethodsItem `json:"extendedResponse"`
	} `json:"responses"`
	Errors []*schemaMethodsErrors
}

func (s schemaMethodsMethods) GetResponse() IMethodItem {
	return s.Responses.Response
}

func (s schemaMethodsMethods) GetExtResponse() IMethodItem {
	return s.Responses.ExtResponse
}

func (s schemaMethodsMethods) GetParameters() []IMethodItem {
	mi := make([]IMethodItem, 0)

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

func (s schemaMethodsItem) GetGoType() string {
	if len(s.Ref) > 0 {
		return getObjectTypeName(s.Ref)
	}

	if fmt.Sprint(s.Type) == SCHEMA_TYPE_ARRAY {
		return s.Items.GetGoType()
	}

	return ""

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
