package main

const (
	OUTPUT_DIR_NAME = "output"

	RESP_DIR_NAME         = "responses"
	RESP_HEADER_TMPL_NAME = "templates/responses.header.template"
	RESP_TMPL_NAME        = "templates/responses.template"

	OBJ_DIR_NAME         = "objects"
	OBJ_HEADER_TMPL_NAME = "templates/objects.header.template"
	OBJ_TMPL_NAME        = "templates/object.template"
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
)

// Response types (string option)
const (
	R_TYPE_INT       string = "integer"
	R_TYPE_STRING    string = "string"
	R_TYPE_ARRAY     string = "array"
	R_TYPE_OBJECT    string = "object"
	R_TYPE_BUILTIN   string = "builtin"
	R_TYPE_BOOLEAN   string = "boolean"
	R_TYPE_INTERFACE string = "interface"
)
