package main

type schemaTyper interface {
	GetGoType() string
	GetDescription() string
}

type schemaTypeChecker interface {
	IsString() bool
	IsInt() bool
	IsBuiltin() bool
	IsArray() bool
	IsObject() bool
	IsBoolean() bool
	IsInterface() bool
}

type schemaTyperChecker interface {
	schemaTyper
	schemaTypeChecker
}

type schemaProperty interface {
	GetProperties() map[string]propertyWrapper
}
