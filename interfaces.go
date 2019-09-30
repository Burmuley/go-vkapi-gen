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

type IMethod interface {
	GetResponse() IMethodItem
	GetExtResponse() IMethodItem
	GetParameters() []IMethodItem
	GetName() string
	IsExtended() bool
}

type IMethodItem interface {
	GetGoType() string
	IsRequired() bool
	GetType() string
	GetName() string
}
