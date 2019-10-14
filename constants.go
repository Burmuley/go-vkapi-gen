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
	SCHEMA_TYPE_INT       string = "integer"
	SCHEMA_TYPE_STRING    string = "string"
	SCHEMA_TYPE_ARRAY     string = "array"
	SCHEMA_TYPE_OBJECT    string = "object"
	SCHEMA_TYPE_BUILTIN   string = "builtin"
	SCHEMA_TYPE_BOOLEAN   string = "boolean"
	SCHEMA_TYPE_INTERFACE string = "interface"
	SCHEMA_TYPE_NUMBER    string = "number"
	SCHEMA_TYPE_UNKNOWN   string = "UNKNOWN"
	SCHEMA_TYPE_MULTIPLE  string = "multiple"
)
