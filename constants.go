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

var (
	outputDirName string = "output"
)

const (
	RESP_DIR_NAME         = "responses"
	RESP_HEADER_TMPL_NAME = "templates/responses.header.template"
	RESP_TMPL_NAME        = "templates/responses.template"

	OBJ_DIR_NAME         = "objects"
	OBJ_HEADER_TMPL_NAME = "templates/objects.header.template"
	OBJ_TMPL_NAME        = "templates/objects.template"

	METHODS_HEADER_TMPL_NAME = "templates/methods.header.template"
	METHODS_TMPL_NAME        = "templates/methods.template"
)

// Response types (string option)
const (
	schemaTypeInt       string = "integer"
	schemaTypeString    string = "string"
	schemaTypeArray     string = "array"
	schemaTypeObject    string = "object"
	schemaTypeBuiltin   string = "builtin"
	schemaTypeBoolean   string = "boolean"
	schemaTypeInterface string = "interface"
	schemaTypeNumber    string = "number"
	schemaTypeUnknown   string = "UNKNOWN"
	schemaTypeMultiple  string = "multiple"
)
