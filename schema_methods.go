package main

type IMethod interface {
	GetResponse() IMethodResponse
	GetExtResponse() IMethodResponse
	IsExtended() bool
	GetParameters() []IMethodParameter
}

type IMethodParameter interface {
	IsRequired() bool
	GetType() string
	GetGoType() string
	GetName() string
}

type IMethodResponse interface {
	GetResponse() string
	GetExtendedResponse() string
}

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
	Name         string                    `json:"name"`
	AccessTokens []string                  `json:"access_token_type"`
	Params       []*schemaMethodsParameter `json:"parameters"`
	Responses    struct {
		Response    *schemaMethodsResponse `json:"response"`
		ExtResponse *schemaMethodsResponse `json:"extendedResponse"`
	} `json:"responses"`
	Errors []*schemaMethodsErrors
}

type schemaMethodsParameter struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Required  bool          `json:"required"`
	Enum      []interface{} `json:"enum"`
	EnumNames []string      `json:"enumNames"`
}

type schemaMethodsResponse struct {
	Ref string `json:"$ref"`
}
