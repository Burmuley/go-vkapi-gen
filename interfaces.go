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

import "text/template"

type IType interface {
	GetGoType() string
	GetDescription() string
	GetType() string
}

type ITypeProperty interface {
	GetProperties(stripPrefix bool) map[string]schemaJSONProperty
}

type IMethod interface {
	GetResponse() IMethodItem
	GetExtResponse() IMethodItem
	GetResponses() []IMethodItem
	GetParameters() []IMethodItem
	GetName() string
	GetDescription() string
	IsExtended() bool
}

type IMethodItem interface {
	IsRequired() bool
	GetName() string
	IType
}

type IGenerator interface {
	Parse(fPath string) error
	Generate(outputDir string) error
}

type IRender interface {
	Render(tmpl *template.Template) ([]byte, error)
}

type IIterator interface {
	Next() (IRender, bool)
	GetKey() string
}
