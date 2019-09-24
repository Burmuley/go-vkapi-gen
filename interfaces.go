package main

type schemaTyper interface {
	GetGoType(stripPrefix bool) []string
	GetDescription() string
	GetProperties(stripPrefix bool) map[string]schemaJSONProperty
}

type schemaTypeChecker interface {
	IsString() bool
	IsInt() bool
	IsBuiltin() bool
	IsArray() bool
	IsObject() bool
	IsBoolean() bool
	IsInterface() bool
	IsNumber() bool
	IsMultiple() bool
}

type schemaTyperChecker interface {
	schemaTyper
	schemaTypeChecker
}
