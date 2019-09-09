package main

const (
	OUTPUT_DIR_NAME = "output"

	RESP_DIR_NAME         = "responses"
	RESP_HEADER_TMPL_NAME = "templates/responses.header.template"
	RESP_TMPL_NAME        = "templates/responses.template"

	OBJ_DIR_NAME         = "objects"
	OBJ_HEADER_TMPL_NAME = "templates/objects.header.template"
	OBJ_TMPL_NAME        = "templates/objects.template"
)

// Responses type constants
const (
	TYPE_INT int = iota
	TYPE_STRING
	TYPE_BUILTIN
	TYPE_ARRAY
	TYPE_OBJECT
	TYPE_BOOLEAN
	TYPE_INTERFACE
	TYPE_NUMBER
)

// Response types (string option)
const (
	SCHEMA_TYPE_INT       string = "integer"
	SCHEMA_TYPE_STRING    string = "string"
	SCHEMA_TYPE_ARRAY     string = "array"
	SCHEMA_TYPE_OBJECT    string = "object"
	SCHEMA_TYPE_BUILTIN   string = "builtin"
	SCHEMA_TYPE_BOOLEAN   string = "boolean"
	SCHEMA_TYPE_INTERFACE string = "interface"
	SCHEMA_TYPE_NUMBER    string = "number"
	SCHEMA_TYPE_UNKNOWN   string = "UNKNOWN"
)

//const ()
