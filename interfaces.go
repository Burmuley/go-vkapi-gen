package main

type schemaTyper interface {
	GetGoType() []string
	GetDescription() string
	GetProperties() map[string]schemaJSONProperty
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
}

type schemaTyperChecker interface {
	schemaTyper
	schemaTypeChecker
}
