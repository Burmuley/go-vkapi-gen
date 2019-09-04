package main

type responseTyper interface {
	GetGoType() string
	GetDescription() string
}

type responseTypeChecker interface {
	IsString() bool
	IsInt() bool
	IsBuiltin() bool
	IsArray() bool
	IsObject() bool
	IsBoolean() bool
	IsInterface() bool
}
