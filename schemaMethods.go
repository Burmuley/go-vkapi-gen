package main

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
	Name         string                `json:"name"`
	AccessTokens []string              `json:"access_token_type"`
	Params       []*schemaJSONProperty `json:"parameters"`
	Resps        []*schemaJSONProperty `json:"responses"`
}
